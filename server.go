package tictactoe

import (
	"errors"
)

type PlayerID string
type BoardID string

// TODO: make interface?
type Server struct {
	games map[PlayerID]map[BoardID]*Board
}

func (s *Server) GetGames(P PlayerID) map[BoardID]*Board {
	return s.games[P]
}

func (s *Server) MakeMove(player PlayerID, board BoardID, move Move) error {
	if _, playerExists := s.games[player]; !playerExists {
		return errors.New("Invalid player")
	}

	targetBoard, boardExists := s.games[player][board];
	if !boardExists {
		return errors.New("Invalid board")
	}

	if valid, err := Valid(&move, targetBoard); !valid {
		return errors.New("Invalid move: " + err.Error())
	}

	return nil;
}
