package test

import "testing"
import . "paperclips/paperclips"

type GameAdapter interface {
	BoardState() Board
	RunMove(m *Move, p PlayerID, tc TurnCount) (*BoardMessage, error)
}

func TestGamePlay(t *testing.T,
	AdapterFactory func(players []PlayerID, startCount int) GameAdapter) {

	TestMoveSequence := func(players []PlayerID, startCount int,
		moves []Move, expectedTurnSequence []PlayerID, expectedWinner PlayerID) {
		Adapter := AdapterFactory(players, startCount)
		// TODO: wrap this?
		currentBoard := Adapter.FirstUpdate()

		Play := func(idx int, m *Move) {
			prevPlayer := currentBoard.WhoseTurn
			msg, err := Adapter.RunMove(m, currentBoard.WhoseTurn, currentBoard.TurnCount)
			if err != nil {
				t.Error("Failed to run move:", err)
			}
			currentBoard = *msg
			if currentBoard.WhoseTurn == prevPlayer {
				t.Error("Failed to advance player counter?!")
			}
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
			Play(i, &m)
			turnSequence = append(turnSequence, currentBoard.WhoseTurn)
		}

		TestTurnSequence(expectedTurnSequence, turnSequence)

		if expectedWinner != "" && expectedWinner != *currentBoard.Winner {
			t.Error("Expected", expectedWinner, "to win, but", *currentBoard.Winner, "won")
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
