package tictactoe

import (
	"errors"
	"strings"
)

type Tile int

const (
	X     Tile = iota
	O     Tile = iota
	Blank Tile = iota
)

func (t *Tile) Render() string {
	switch *t {
	case X:
		return "X"
	case O:
		return "O"
	case Blank:
		return " "
	default:
		panic("Unknown tile type?!")
	}
}

type Board [3][3]Tile

func (m *Board) Render() string {
	lines := make([]string, 0)
	for x := 0; x < 3; x++ {
		line := make([]string, 0)
		for y := 0; y < 3; y++ {
			line = append(line, m[x][y].Render())
		}
		lines = append(lines, strings.Join(line, "|"))
	}
	return strings.Join(lines, "-----")
}

type Move struct {
	x, y int
	tile Tile
}

func (m *Move) Valid() (bool, error) {
	if (m.x < 0 || m.x > 2 ||
		m.y < 0 || m.y > 2) {
		return false, errors.New("Out of bounds")
	}
	if (m.tile != X && m.tile != O) {
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

