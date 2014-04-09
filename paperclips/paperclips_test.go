package paperclips

import "testing"

func TestRender(t *testing.T) {
	var data Board
	data.PaperclipCount = 6
	if data.Render() != "6" {
		t.Errorf("Failed to properly render board!")
	}
}

func TestMoveBounds(t *testing.T) {
	ExpectValid := func(m Move) {
		if val, err := m.Valid(); !val || !(err == nil) {
			t.Error(err)
			t.Error("Move", m, "was supposed to be valid!")
		}
	}
	ExpectInvalid := func(m Move) {
		if val, err := m.Valid(); val || (err == nil) {
			t.Error(err)
			t.Error("Move", m, "was supposed to be invalid!")
		}
	}

	ExpectValid(TakeOne)
	ExpectValid(TakeTwo)
	ExpectInvalid(3)
	ExpectInvalid(0)
}

func TestGamePlay(t *testing.T) {
	TestMoveSequence := func(players []PlayerID, startCount int, moves []Move, expectedTurnSequence []PlayerID, expectedWinner PlayerID) {
		board := NewBoard(players, startCount, "TestGame")
		Play := func(idx int, m *Move) {
			prevPlayer := board.CurrentPlayer()
			if err := board.Apply(m); err != nil {
				t.Error("Failed to apply move", idx, ":", m, ":", err.Error())
			}
			if board.CurrentPlayer() == prevPlayer {
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
			turnSequence = append(turnSequence, board.CurrentPlayer())
		}

		TestTurnSequence(expectedTurnSequence, turnSequence)

		if expectedWinner != "" && expectedWinner != board.WinningPlayer() {
			t.Error("Expected", expectedWinner, "to win, but", board.WinningPlayer(), "won")
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
