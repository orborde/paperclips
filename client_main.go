package main

import (
	"fmt"
	"log"
	"net"
	"paperclips/paperclips"
	"strconv"
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
	defer client.Close()

	var name paperclips.PlayerID
	for i := 0; ; i++ {
		name = paperclips.PlayerID("Player" + strconv.Itoa(i))
		if client.RegisterPlayer(name) == nil {
			break
		}
	}
	log.Println("My name is", name)

	if names, err := client.GetPlayerList(); err != nil {
		log.Fatal("Failed to fetch player list:", err)
	} else {
		log.Println("Players online are:", names)
	}
}
