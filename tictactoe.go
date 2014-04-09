package tictactoe

import (
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

func (m *Move) Valid() bool {
	if m.x < 0 || m.x > 2 ||
		m.y < 0 || m.y > 2 ||
		(m.tile != X && m.tile != O) {
		return false
	}
	return true
}

func Valid(m *Move, b *Board) {
	if !m.Valid() ||
		b[m.x][m.y] != Blank {
		return false
	}
	return true
}

