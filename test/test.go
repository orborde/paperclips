package test

import "testing"
import . "paperclips/paperclips"

type GameAdapter interface {
	BoardState() *Board
	RunMove(m *Move, p PlayerID) (*Board, error)
}

func TestGamePlay(t *testing.T,
	AdapterFactory func(players []PlayerID, startCount int) GameAdapter) {

	TestMoveSequence := func(players []PlayerID, startCount int,
		moves []Move, expectedTurnSequence []PlayerID, expectedWinner PlayerID) {
		Adapter := AdapterFactory(players, startCount)

		Play := func(idx int, m *Move) bool {
			prevPlayer := Adapter.BoardState().CurrentPlayer()
			_, err := Adapter.RunMove(m, Adapter.BoardState().CurrentPlayer())
			if err != nil {
				t.Error("Failed to run move:", err)
				return false
			}

			if Adapter.BoardState().CurrentPlayer() == prevPlayer {
				t.Error("Failed to advance player counter?!")
				return false
			}
			return true
		}

		TestTurnSequence := func(expected, actual []PlayerID) {
			ok := len(expected) == len(actual)
			for i := range expected {
				if expected[i] != actual[i] {
					ok = false
				}
			}
			if !ok {
				t.Error("Unexpected turn sequence", actual, ", expected", expected)
			}
		}

		turnSequence := []PlayerID{players[0]}
		for i, m := range moves {
			if !Play(i, &m) {
				return
			}
			turnSequence = append(turnSequence, Adapter.BoardState().CurrentPlayer())
		}

		TestTurnSequence(expectedTurnSequence, turnSequence)

		if expectedWinner != "" && expectedWinner != Adapter.BoardState().WinningPlayer {
			t.Error("Expected", expectedWinner, "to win, but", Adapter.BoardState().WinningPlayer, "won")
		}
	}

	TestMoveSequence([]PlayerID{"a", "b"}, 5,
		[]Move{1, 2, 1, 1}, []PlayerID{"a", "b", "a", "b", "a"}, "b")

	TestMoveSequence([]PlayerID{"a", "b", "c"}, 5,
		[]Move{1, 2, 1, 1}, []PlayerID{"a", "b", "c", "a", "b"}, "a")

	TestMoveSequence([]PlayerID{"b", "a"}, 4,
		[]Move{1, 1, 2}, []PlayerID{"b", "a", "b", "a"}, "b")

	TestMoveSequence([]PlayerID{"b", "a"}, 4,
		[]Move{1, 1, 1, 1}, []PlayerID{"b", "a", "b", "a", "b"}, "a")
}
