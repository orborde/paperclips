package paperclips

import (
	"testing"

	"paperclips/test"
)

type RawGameAdapter struct {
	paperclips.GameAdapter
	game        Game
	firstUpdate *BoardMessage
}

func NewRawGameAdapter() *RawGameAdapter {
	ret := &RawGameAdapter{NewGame(players, PaperclipCount(startCount)), nil}
	firstUpdate <- ret.game.FirstUpdate
	return ret
}

func (a *RawGameAdapter) FirstUpdate() BoardMessage {
	return *a.firstUpdate
}

func (a *RawGameAdapter) RunMove(m *Move, p PlayerID, tc TurnCount) (*BoardMessage, error) {
	result := make(chan MoveResult)
	a.game.Moves <- MoveMessage{*m, player, turnCount, result}
	msg := <-result
	return msg.BoardMessage, msg.error
}

func TestRender(t *testing.T) {
	var data Board
	data.PaperclipCount = 6
	if data.Render() != "6" {
		t.Errorf("Failed to properly render board!")
	}
}

func TestMoveBounds(t *testing.T) {
	ExpectValid := func(m Move) {
		if val, err := m.Valid(); !val || !(err == nil) {
			t.Error(err)
			t.Error("Move", m, "was supposed to be valid!")
		}
	}
	ExpectInvalid := func(m Move) {
		if val, err := m.Valid(); val || (err == nil) {
			t.Error(err)
			t.Error("Move", m, "was supposed to be invalid!")
		}
	}

	ExpectValid(TakeOne)
	ExpectValid(TakeTwo)
	ExpectInvalid(3)
	ExpectInvalid(0)
}
