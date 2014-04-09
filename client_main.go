package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"paperclips/paperclips"
)

const host string = "localhost"
const port string = "19996"

var address string = host + ":" + port

func main() {
	fmt.Println("Starting up RPC server...")
	var server paperclips.RPCServer
	rpc.Register(&server)
	fmt.Println("Now listening on", address)
	l, e := net.Listen("tcp", address)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	rpc.Accept(l)
}
