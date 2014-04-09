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

type Board struct {
	PaperclipCount int
	Players        []string
	// Indices into the Players[] array
	CurrentPlayerIndex int
	WinningPlayerIndex int
	// TODO: Move history?
	BoardID string
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
		b.WinningPlayerIndex = b.CurrentPlayerIndex
	}

	b.CurrentPlayerIndex = (b.CurrentPlayerIndex + 1) % len(b.Players)
	return nil
}

func (b *Board) GameOver() bool {
	return b.PaperclipCount == 0
}

func (b *Board) CurrentPlayer() string {
	return b.Players[b.CurrentPlayerIndex]
}

func (b *Board) WinningPlayer() string {
	if !b.GameOver() {
		return ""
	}
	return b.Players[b.WinningPlayerIndex]
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
	if int(*m) > b.PaperclipCount {
		return false, errors.New("Not enough paperclips")
	}
	return true, nil
}
