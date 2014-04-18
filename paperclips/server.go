package paperclips

import (
	"errors"
	"fmt"
	"strconv"
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
type PaperclipServer interface {
	// Returns a list of games the player identified by PlayerID is
	// currently participating in on the server.
	GetGames(PlayerID) (map[BoardID]*Board, error)
	// Makes Move on Board, updating the server's game state.
	MakeMove(Player PlayerID, Board BoardID, Move Move) error
	// Registers a new Player as participating in the game server.
	NewPlayer(Name PlayerID) error
	// Creates a new game on the server between the listed Players.
	NewGame(Players []PlayerID, StartCount int) (BoardID, error)
	// Grabs the player list.
	GetPlayerList() ([]PlayerID, error)
}

// Clients will poll the GetGames interface periodically to receive a
// list of active Boards; later, they will send MakeMove() RPCs back
// to submit the user's moves on her Boards. The RPC interface is
// designed to be friendly to a client that communicates with the
// server entirely in the background so that the user never has to
// wait for network round-trips while making moves; the client is
// expected to store the downloaded list of active games and to queue
// up the moves made for later delivery via MakeMove.

type Server struct {
	games       map[PlayerID]map[BoardID]*Board
	nextBoardId uint64
}

func NewServer() *Server {
	return &Server{make(map[PlayerID]map[BoardID]*Board), 0}
}

func (s *Server) PlayerExists(P PlayerID) bool {
	_, ret := s.games[P]
	return ret
}

func (s *Server) NewPlayer(Name PlayerID) error {
	if s.PlayerExists(Name) {
		return errors.New("Player " + string(Name) + " already exists on server")
	}
	s.games[Name] = make(map[BoardID]*Board)
	return nil
}

func (s *Server) GetPlayerList() ([]PlayerID, error) {
	ret := make([]PlayerID, 0)
	for p := range s.games {
		ret = append(ret, p)
	}
	return ret, nil
}

func (s *Server) getNextBoardId() BoardID {
	ret := BoardID(strconv.FormatUint(s.nextBoardId, 10))
	s.nextBoardId++
	return ret
}

func (s *Server) NewGame(Players []PlayerID, StartCount int) (BoardID, error) {
	for _, p := range Players {
		if !s.PlayerExists(p) {
			//return errors.New("Player " + string(p) + " does not exist on server")
			if err := s.NewPlayer(p); err != nil {
				return "", err
			}
		}
	}

	board := NewBoard(Players, StartCount)
	ID := s.getNextBoardId()
	for _, p := range Players {
		s.games[p][ID] = board
	}
	return ID, nil
}

func (s *Server) GetGames(P PlayerID) (map[BoardID]*Board, error) {
	if !s.PlayerExists(P) {
		return nil, errors.New(fmt.Sprint("Player", P, "does not exist."))
	}
	return s.games[P], nil
}

func (s *Server) MakeMove(player PlayerID, board BoardID, move Move) error {
	if !s.PlayerExists(player) {
		return errors.New("Invalid player")
	}

	targetBoard, boardExists := s.games[player][board]
	if !boardExists {
		return errors.New("Invalid board")
	}

	if valid, err := Valid(&move, targetBoard); !valid {
		return errors.New("Invalid move: " + err.Error())
	}

	return targetBoard.Apply(&move)
}
