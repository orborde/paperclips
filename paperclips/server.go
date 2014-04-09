package paperclips

import (
	"errors"
	"fmt"
	"strconv"
)

// TODO: make interface?
type Server struct {
	games       map[PlayerID]map[BoardID]*Board
	nextBoardId uint64
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

	board := NewBoard(Players, StartCount, s.getNextBoardId())
	ID := board.ID
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

	return nil
}

type RPCServer struct {
	server Server
}

func (s *RPCServer) RegisterPlayer(P PlayerID, _ *struct{}) error {
	return s.server.NewPlayer(P)
}

type RPCNewGameArgs struct {
	Players    []PlayerID
	StartCount int
}

func (s *RPCServer) NewGame(Args RPCNewGameArgs, Id *BoardID) error {
	id, err := s.server.NewGame(Args.Players, Args.StartCount)
	*Id = id
	return err
}

func (s *RPCServer) GetGames(P PlayerID, Ret *map[BoardID]*Board) error {
	var err error
	*Ret, err = s.server.GetGames(P)
	return err
}

type RPCMove struct {
	Player PlayerID
	// TODO: check that player is on board
	// TODO: handle moves from stale board states
	BoardID BoardID
	Move    Move
}

func (s *RPCServer) MakeMove(Args RPCMove, _ *struct{}) error {
	return s.server.MakeMove(Args.Player, Args.BoardID, Args.Move)
}
