package main

import (
	"fmt"
	"log"
	"net"
	"paperclips/paperclips"
)

const host string = "localhost"
const port string = "19996"

var address string = host + ":" + port

func main() {
	fmt.Println("Connecting to expected RPC server at", address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal("connect error: ", err)
	}

	client := paperclips.NewRPCClient(conn)
	{
		err := client.RegisterPlayer("potato-head")
		if err != nil {
			log.Fatal("call error: ", err)
		}
	}
}
