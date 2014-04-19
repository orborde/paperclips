package test

import (
	"log"
	"net"
	"net/rpc"
	"testing"
)

import . "paperclips/paperclips"

const address string = "localhost:34824"

func NewRPCServerGameAdapter(players []PlayerID, startCount int) GameAdapter {
	rpcServer := NewRPCServer()
	rpc.Register(rpcServer)
	l, e := net.Listen("tcp", address)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go rpc.Accept(l)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal("connect error: ", err)
	}

	rpcClient := NewRPCClient(conn)

	ret := &ServerGameAdapter{server: rpcClient}
	return ret
}

func TestRPCServer(t *testing.T) {
	TestGamePlay(t, NewRPCServerGameAdapter)
}
