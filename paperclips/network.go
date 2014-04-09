package paperclips

import (
	"log"
)

type RPCServer struct {
	server *Server
}

func NewRPCServer() *RPCServer {
	return &RPCServer{NewServer()}
}

func (s *RPCServer) RegisterPlayer(P PlayerID, _ *struct{}) error {
	log.Println("Registering player", P)
	return s.server.NewPlayer(P)
}

type RPCNewGameArgs struct {
	Players    []PlayerID
	StartCount int
}

func (s *RPCServer) NewGame(Args RPCNewGameArgs, Id *BoardID) error {
	log.Println("Configuring a new game with args", Args)
	id, err := s.server.NewGame(Args.Players, Args.StartCount)
	*Id = id
	return err
}

func (s *RPCServer) GetGames(P PlayerID, Ret *map[BoardID]*Board) error {
	log.Println("Fetching games for player", P)
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
	log.Println("Processing move", Args)
	return s.server.MakeMove(Args.Player, Args.BoardID, Args.Move)
}
