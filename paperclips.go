package paperclips

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Move int

const (
	TakeOne Move = 1
	TakeTwo Move = 2
)

type Board struct {
	PaperclipCount int
	Players []string
	// Indices into the Players[] array
	NextPlayer int
	WinningPlayer int
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
	b[move.x][move.y] = move.tile
	return nil
}

func (b *Board) Winner() Tile {
	checkWin := func(x, y []int) Tile {

	}
}

type Move struct {
	x, y int
	tile Tile
}

func (m *Move) Valid() (bool, error) {
	if m.x < 0 || m.x > 2 ||
		m.y < 0 || m.y > 2 {
		return false, errors.New("Out of bounds")
	}
	if m.tile != X && m.tile != O {
		return false, errors.New("Not a placeable tile")
	}
	return true, nil
}

func Valid(m *Move, b *Board) (bool, error) {
	if valid, err := m.Valid(); !valid {
		return false, err
	}
	if b[m.x][m.y] != Blank {
		return false, errors.New("Tile is occupied")
	}
	return true, nil
}
