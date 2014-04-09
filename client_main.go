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
	defer client.Close()

	players := []paperclips.PlayerID{"kim", "joe"}
	for _, p := range players {
		err := client.RegisterPlayer(p)
		if err != nil {
			log.Println("call error: ", err)
		}
	}

	{
		boardId, err := client.NewGame(players, 5)
		if err != nil {
			log.Fatal("Failed to create a game")
		}
		fmt.Println("Game created:", boardId)
	}

	for _, p := range players {
		data, err := client.GetGames(p)
		if err != nil {
			log.Fatal("Fetch failure:", err)
		}
		log.Println(p, "games are", data)
	}
}
