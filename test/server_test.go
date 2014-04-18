package test

import (
	"fmt"
	"log"
	"testing"
)

import . "paperclips/paperclips"

type ServerGameAdapter struct {
	server PaperclipServer
	board  BoardID
}

func NewLocalServerGameAdapter(players []PlayerID, startCount int) GameAdapter {
	return NewServerGameAdapter(NewServer(), players, startCount)
}

func NewServerGameAdapter(server PaperclipServer, players []PlayerID, startCount int) GameAdapter {
	ret := &ServerGameAdapter{server: server}

	// Set up players on server
	for _, p := range players {
		if err := ret.server.NewPlayer(p); err != nil {
			log.Fatal("Failed to create players:", err)
		}
	}

	// Create a game among the players
	{
		var err error
		ret.board, err = ret.server.NewGame(players, startCount)
		if err != nil {
			log.Fatal(err)
		}
	}

	return ret
}

func (a *ServerGameAdapter) BoardState() (board *Board) {
	// Fetch the players list.
	players := a.server.GetPlayerList()

	// Grab a player off it and fetch his games.
	p := players[0]
	games, err := a.server.GetGames(p)
	if err != nil {
		panic(fmt.Sprint("Server failed to give back games:", err))
	}

	// It should be the only game in the list.
	return games[a.board]
}

func (a *ServerGameAdapter) RunMove(m *Move, p PlayerID) (*Board, error) {
	err := a.server.MakeMove(p, a.board, *m)
	return a.BoardState(), err
}

func TestServer(t *testing.T) {
	TestGamePlay(t, NewLocalServerGameAdapter)
}
