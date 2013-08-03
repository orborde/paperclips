package tictactoe

import (
	"strings"
)

type Tile int

const (
	X tile = iota
	O tile = iota
	Blank tile = iota
)

func (t *Tile) Render() string {
	switch t {
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
	lines := [](string*)
	for x = 0; x < 3; x++ {
		line := [](string*)
		for y = 0; y < 3; y++ {
			append(line, &m[x][y].Render())
		}
		append(lines, strings.Join(line, "|"))
	}
	return strings.Join(line, "-----")
}
