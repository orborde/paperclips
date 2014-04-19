package test

import (
	"log"
	"net"
	"net/rpc"
	"testing"
)

import . "paperclips/paperclips"

const address string = "localhost:0"

func NewRPCServerGameAdapter(players []PlayerID, startCount int) GameAdapter {
	rpcServer := rpc.NewServer()
	gameServer := NewRPCServer()
	rpcServer.Register(gameServer)
	l, e := net.Listen("tcp", address)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go rpcServer.Accept(l)

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		log.Fatal("connect error: ", err)
	}

	rpcClient := NewRPCClient(conn)

	ret := NewServerGameAdapter(rpcClient, players, startCount)
	return ret
}

func TestRPCServer(t *testing.T) {
	TestGamePlay(t, NewRPCServerGameAdapter)
}
