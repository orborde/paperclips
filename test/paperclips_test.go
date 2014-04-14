package test

import "testing"

import . "paperclips/paperclips"

type RawGameAdapter struct {
	game Board
}

func NewRawGameAdapter(players []PlayerID, startCount int) GameAdapter {
	return &RawGameAdapter{*NewBoard(players, startCount)}
}

func (a *RawGameAdapter) BoardState() *Board {
	return &a.game
}

func (a *RawGameAdapter) RunMove(m *Move, p PlayerID) (*Board, error) {
	// TODO: What about wrong player?
	err := a.game.Apply(m)
	return &a.game, err
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

func TestRaw(t *testing.T) {
	TestGamePlay(t, NewRawGameAdapter)
}
