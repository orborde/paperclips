package ai

import (
	"log"

	. "paperclips/paperclips"
)

func ValidMoves(b *Board) []Move {
	ret := make([]Move, 0)
	for _, m := range Moves {
		if valid, _ := Valid(&m, b); valid {
			ret = append(ret, m)
		}
	}
	return ret
}

type Progression struct {
	Board *Board
	Move  *Move
}

func Progressions(b *Board) []Progression {
	ret := make([]Progression, 0)
	for _, m := range ValidMoves(b) {
		nb := *b
		nb.Apply(&m)
		ret = append(ret, Progression{&nb, &m})
	}
	return ret
}

// Stupid go glue
func setToList(m map[PlayerID]bool) []PlayerID {
	ret := make([]PlayerID, 0)
	for k := range m {
		log.Print("berkle: ", k)
		ret = append(ret, k)
	}
	return ret
}

// Returns "" for "indeterminate"
func ComputeWinner(b *Board) PlayerID {
	if b.GameOver() {
		return b.WinningPlayer
	}

	// Figure out what the possible winners are in the subtree.
	possibleWinners := make(map[PlayerID]bool)
	for _, p := range Progressions(b) {
		subWinner := ComputeWinner(p.Board)
		possibleWinners[subWinner] = true
	}

	possibleList := setToList(possibleWinners)
	log.Print(b.PaperclipCount, b.CurrentPlayer(), possibleList)

	// Are we one of the possible winners? We are the CHAMPIONS.
	if _, ok := possibleWinners[b.CurrentPlayer()]; ok {
		log.Print("I'm the winrar!")
		return b.CurrentPlayer()
	}

	// Is there only one possible winner? Then that's the winner for this subtree.
	if len(possibleWinners) == 1 {
		log.Print("There can be only one: ", len(possibleList))
		return possibleList[0]
	}

	// a winner is not me
	log.Print("A winner is nobody")
	return ""
}
