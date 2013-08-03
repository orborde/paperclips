package tictactoe

import (
	"strings"
)

type Tile int

const (
	X Tile = iota
	O Tile = iota
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

type Map [3][3]Tile

func (m *Map) Render() string {
	lines := make([]*string, 0)
	for x := 0; x < 3; x++ {
		line := make([]*string, 0)
		for y := 0; y < 3; y++ {
			append(line, m[x][y].Render())
		}
		append(lines, strings.Join(line, "|"))
	}
	return strings.Join(line, "-----")
}
