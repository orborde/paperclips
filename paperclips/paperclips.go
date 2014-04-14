package paperclips

import (
	"errors"
	"fmt"
	"strconv"
)

type Move int

const (
	TakeOne Move = 1
	TakeTwo Move = 2
)

type PlayerID string
type BoardID string

type Board struct {
	PaperclipCount     int
	Players            []PlayerID
	currentPlayerIndex int
	WinningPlayer      Player
	// TODO: Move history?
	ID BoardID
}

func NewBoard(Players []PlayerID, StartCount int, Id BoardID) *Board {
	return &Board{PaperclipCount: StartCount, Players: Players, ID: Id}
}

func (b *Board) Render() string {
	return strconv.Itoa(b.PaperclipCount)
}

func (b *Board) Apply(move *Move) error {
	if valid, err := Valid(move, b); !valid {
		return errors.New(fmt.Sprint("Could not apply move ", move, ":", err.Error()))
	}
	b.PaperclipCount -= int(*move)

	if b.GameOver() {
		b.WinningPlayer = b.CurrentPlayer()
	}

	b.currentPlayerIndex = (b.currentPlayerIndex + 1) % len(b.Players)
	return nil
}

func (b *Board) GameOver() bool {
	return b.PaperclipCount == 0
}

func (b *Board) CurrentPlayer() PlayerID {
	return b.Players[b.currentPlayerIndex]
}

func (b *Board) WinningPlayer() PlayerID {
	if !b.GameOver() {
		return ""
	}
	return b.WinningPlayer
}

func (m *Move) Valid() (bool, error) {
	if *m != TakeOne && *m != TakeTwo {
		return false, errors.New("Not a valid move type")
	}
	return true, nil
}

func Valid(m *Move, b *Board) (bool, error) {
	if valid, err := m.Valid(); !valid {
		return false, err
	}
	if b.GameOver() {
		return false, errors.New("Game has already ended")
	}
	if int(*m) > b.PaperclipCount {
		return false, errors.New("Not enough paperclips")
	}
	return true, nil
}
