package tictactoe

type PlayerID string
type BoardID string

// TODO: make interface?
type Server struct {
	games map[PlayerID]map[BoardID]Map
}

func (s *Server) GetGames(P PlayerID) map[BoardID]Map {
	return s.games[P]
}

func (s *Server) MakeMove(player PlayerID, board BoardID, move Move) error {
	if _, playerExists := s.games[player]; !playerExists {
		return errors.New("Invalid player")
	}

	var targetBoard *Map
	var boardExists bool
	targetBoard, boardExists = s.games[player][board];
	if !boardExists {
		return errors.New("Invalid board")
	}

	if !move.Valid(board) {
		return errors.New("Invalid move")
	}

	return nil;
}
