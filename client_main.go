package main

import (
	"fmt"
	"log"
	//"math/rand"
	"net"
	"paperclips/paperclips"
	"strconv"
	"time"
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

	for {
		log.Println("Fetching my games...")
		games, err := client.GetGames(name)
		if err != nil {
			log.Fatal("Failed to fetch games:", err)
		}

		// TODO: Don't have the server dump all the games ever.
		activeGames := make([]paperclips.BoardID, 0)
		for id, game := range games {
			if !game.GameOver() {
				activeGames = append(activeGames, id)
			}
		}

		log.Println("Fetched", len(games), "games, of which", len(activeGames), "are active.")

		log.Println("Time to sleep")
		time.Sleep(10 * time.Second)
	}
}
