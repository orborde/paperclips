package paperclips

/*
import (
	"io"
	"log"
	"net/rpc"
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

func (s *RPCServer) GetPlayerList(_ bool, Ret *[]PlayerID) error {
	log.Println("Fetching player list")
	*Ret = s.server.GetPlayerList()
	return nil
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

// Derp derp derp
type RPCClient struct {
	*rpc.Client
}

func NewRPCClient(Conn io.ReadWriteCloser) *RPCClient {
	return &RPCClient{rpc.NewClient(Conn)}
}

func (c *RPCClient) RegisterPlayer(P PlayerID) error {
	return c.Call("RPCServer.RegisterPlayer", P, nil)
}

func (c *RPCClient) GetPlayerList() ([]PlayerID, error) {
	var ret []PlayerID
	err := c.Call("RPCServer.GetPlayerList", false, &ret)
	return ret, err
}

func (c *RPCClient) NewGame(Players []PlayerID, StartCount int) (BoardID, error) {
	var id BoardID
	err := c.Call("RPCServer.NewGame", RPCNewGameArgs{Players, StartCount}, &id)
	return id, err
}

func (c *RPCClient) GetGames(P PlayerID) (map[BoardID]*Board, error) {
	var ret map[BoardID]*Board
	err := c.Call("RPCServer.GetGames", P, &ret)
	return ret, err
}

func (c *RPCClient) MakeMove(Player PlayerID, Board BoardID, Move Move) error {
	return c.Call("RPCServer.MakeMove",
		RPCMove{Player, Board, Move}, nil)
}
*/
