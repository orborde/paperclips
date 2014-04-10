package paperclips

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

// An implementation of a (fairly generic) server for turn-based games
// like Tic-Tac-Toe, or, in this case, a much simpler game I call
// Paperclips.
//
// The server's job is to act as a store of game states. Each game in
// progress is represented by a Board object. A Board contains the
// current game state (board state, whose turn it is, who is playing)
// and is identified by a BoardID (a unique string identifier). There
// are also Move objects, which can be submitted to the server to make
// a move on a Board. When this happens, the Board is updated with the
// new game state.
//
// Clients interact with the server via a polling RPC interface,
// consisting of a couple of major methods:
//
// - GetGames(PlayerID) returns a list of games the player identified
//   by PlayerID is currently participating in on the server.
// - MakeMove(Board, Move) makes Move on Board, updating the server's
//   game state.
// - NewPlayer(PlayerID) registers a new Player as participating in
//   the game server.
// - NewGame(Players []PlayerID, GameOptions) creates a new game on
//   the server between the listed Players.
//
// Clients will poll the GetGames interface periodically to receive a
// list of active Boards; later, they will send MakeMove() RPCs back
// to submit the user's moves on her Boards. The RPC interface is
// designed to be friendly to a client that communicates with the
// server entirely in the background so that the user never has to
// wait for network round-trips while making moves; the client is
// expected to store the downloaded list of active games and to queue
// up the moves made for later delivery via MakeMove.

type Server struct {
	directoryLock sync.RWMutex
	games         map[PlayerID]map[BoardID]*Game
	boards        map[PlayerID]map[BoardID]*BoardMessage
	nextBoardId   uint64
}

func NewServer() *Server {
	return &Server{games: make(map[PlayerID]map[BoardID]*Game)}
}

func (s *Server) PlayerExists(P PlayerID) bool {
	s.directoryLock.RLock()
	defer s.directoryLock.RUnlock()
	_, ret := s.games[P]
	return ret
}

func (s *Server) NewPlayer(Name PlayerID) error {
	s.directoryLock.Lock()
	defer s.directoryLock.Unlock()
	if s.PlayerExists(Name) {
		return errors.New("Player " + string(Name) + " already exists on server")
	}
	s.games[Name] = make(map[BoardID]*Game)
	return nil
}

func (s *Server) GetPlayerList() []PlayerID {
	s.directoryLock.RLock()
	defer s.directoryLock.RUnlock()
	ret := make([]PlayerID, len(s.games))
	for p := range s.games {
		ret = append(ret, p)
	}
	return ret
}

func (s *Server) getNextBoardId() BoardID {
	ret := BoardID(strconv.FormatUint(s.nextBoardId, 10))
	s.nextBoardId++
	return ret
}

func (s *Server) NewGame(Players []PlayerID, StartCount int) (BoardID, error) {
	s.directoryLock.Lock()
	defer s.directoryLock.Unlock()
	for _, p := range Players {
		if !s.PlayerExists(p) {
			//return errors.New("Player " + string(p) + " does not exist on server")
			if err := s.NewPlayer(p); err != nil {
				return "", err
			}
		}
	}

	game := NewGame(Players, PaperclipCount(StartCount))
	ID := s.getNextBoardId()
	for _, p := range Players {
		s.games[p][ID] = game
	}
	return "", nil
}

func (s *Server) GetGames(P PlayerID) (map[BoardID]*BoardMessage, error) {
	s.directoryLock.RLock()
	defer s.directoryLock.RUnlock()
	if !s.PlayerExists(P) {
		return nil, errors.New(fmt.Sprint("Player", P, "does not exist."))
	}
	return s.boards[P], nil
}

func (s *Server) MakeMove(player PlayerID, board BoardID, move Move, turnCount TurnCount) error {
	// Theoretically, we only need to RLock here since MakeMove
	// dispatches a message and doesn't piddle with internal state.
	s.directoryLock.Lock()
	defer s.directoryLock.Unlock()
	if !s.PlayerExists(player) {
		return errors.New("Invalid player")
	}

	targetGame, gameExists := s.games[player][board]
	if !gameExists {
		return errors.New("Invalid game")
	}

	// ...Except we want to update the board cache here, and manual
	// locks is danger zone.
	msg := MoveMessage{move, player, turnCount, make(chan MoveResult)}
	targetGame.Moves <- msg
	// Wait for the coroutine to reply
	result := <-msg.Result
	if result.Error != nil {
		return result.Error
	}

	// Update the board cache
	s.boards[player][board] = result.BoardMessage

	return nil
}
