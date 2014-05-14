package ai

import (
	"testing"

	. "paperclips/paperclips"
)

func TestComputeWinner(t *testing.T) {
	players := []PlayerID{"a", "b"}

	checkWinner := func(startCount int, winner PlayerID) {
		if actual := ComputeWinner(NewBoard(players, startCount)); winner != actual {
			t.Error("Wrong winner:", actual, ", expected", winner)
		}
	}

	checkWinner(1, "a")
	checkWinner(2, "a")
	checkWinner(3, "b")
	checkWinner(4, "a")
}
