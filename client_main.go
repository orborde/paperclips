package main

import (
	"fmt"
	"log"
	"net/rpc"
	"paperclips/paperclips"
)

const host string = "localhost"
const port string = "19996"

var address string = host + ":" + port

func main() {
	fmt.Println("Connecting to expected RPC server at", address)
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		log.Fatal("connect error: ", err)
	}

	args := paperclips.PlayerID("potato-head")
	{
		err := client.Call("RPCServer.RegisterPlayer", args, nil)
		if err != nil {
			log.Fatal("call error: ", err)
		}
	}
}
