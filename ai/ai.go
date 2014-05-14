package ai

import (
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
		ret = append(ret, k)
	}
	return ret
}

// Returns "" for "indeterminate"
func GoComputeWinner(b *Board) chan PlayerID {
	ret := make(chan PlayerID)

	go func() {
		if b.GameOver() {
			ret <- b.WinningPlayer
			return
		}

		// Figure out what the possible winners are in the subtree.
		possibleWinners := make(map[PlayerID]bool)
		var promises []chan PlayerID
		for _, p := range Progressions(b) {
			pr := GoComputeWinner(p.Board)
			promises = append(promises, pr)
		}

		for _, pr := range promises {
			subWinner := <-pr
			possibleWinners[subWinner] = true
		}

		possibleList := setToList(possibleWinners)

		// Are we one of the possible winners? We are the CHAMPIONS.
		if _, ok := possibleWinners[b.CurrentPlayer()]; ok {
			ret <- b.CurrentPlayer()
			return
		}

		// Is there only one possible winner? Then that's the winner for this subtree.
		if len(possibleWinners) == 1 {
			ret <- possibleList[0]
			return
		}

		// a winner is not me
		ret <- ""
		return
	}()

	return ret
}

func ComputeWinner(b *Board) PlayerID {
	return <-GoComputeWinner(b)
}
