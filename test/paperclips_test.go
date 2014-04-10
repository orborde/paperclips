package test

import "testing"

import . "paperclips/paperclips"

type RawGameAdapter struct {
	game        Game
	firstUpdate *BoardMessage
}

func NewRawGameAdapter(players []PlayerID, startCount int) *RawGameAdapter {
	ret := &RawGameAdapter{*NewGame(players, PaperclipCount(startCount)), nil}
	tmp := <-ret.game.FirstUpdate
	ret.firstUpdate = &tmp
	return ret
}

func (a *RawGameAdapter) FirstUpdate() BoardMessage {
	return *a.firstUpdate
}

func (a *RawGameAdapter) RunMove(m *Move, p PlayerID, tc TurnCount) (*BoardMessage, error) {
	result := make(chan MoveResult)
	a.game.Moves <- MoveMessage{*m, p, tc, result}
	msg := <-result
	return msg.BoardMessage, msg.Error
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
