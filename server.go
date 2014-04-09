package paperclips

import (
	"errors"
	"fmt"
)

type PlayerID string
type BoardID string

// TODO: make interface?
type Server struct {
	games map[PlayerID]map[BoardID]*Board
}

func (s *Server) PlayerExists(P PlayerID) bool {
	_, ret := s.games[P]
	return ret
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

func (s *RPCServer) GetGames(P PlayerID, Ret *map[BoardID]*Board) error {
	var err error
	*Ret, err = s.server.GetGames(P)
	return err
}

type RPCMove struct {
	Player PlayerID
	// TODO: check that player is on board
	BoardID BoardID
	Move    Move
}

func (s *RPCServer) MakeMove(Args RPCMove, _ *struct{}) error {
	return s.server.MakeMove(Args.Player, Args.BoardID, Args.Move)
}
