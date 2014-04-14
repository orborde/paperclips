package test

import "testing"

import . "paperclips/paperclips"

import "log"

type ServerGameAdapter struct {
	GameAdapter
	server      Server
	board       BoardID
	firstUpdate *BoardMessage
}

func NewServerGameAdapter(players []PlayerID, startCount int) GameAdapter {
	ret := ServerGameAdapter{server: *NewServer()}

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

func (a ServerGameAdapter) FirstUpdate() BoardMessage {
	if a.firstUpdate == nil {
		log.Fatal("IS NIL")
	}
	return *(a.firstUpdate)
}

func (a ServerGameAdapter) RunMove(m *Move, p PlayerID, tc TurnCount) (*BoardMessage, error) {
	result := make(chan MoveResult)
	a.game.Moves <- MoveMessage{*m, p, tc, result}
	msg := <-result
	return msg.BoardMessage, msg.Error
}

func (a ServerGameAdapter) GetGame() (ID BoardID) {
	// Fetch the players list.
	players := a.server.GetPlayerList()

	// Grab a player off it and fetch his games.
	p := players[0]
	games, err := a.server.GetGames(p)
	if err != nil {
		panic("Server failed to give back games!")
	}

	// It should be the only game in the list.
	for _, boardID := range games {
		ID = boardID
		break
	}
	return ID
}

func TestServer(t *testing.T) {
	TestGamePlay(t, NewServerGameAdapter)
}
