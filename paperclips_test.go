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
	TestMoveSequence := func(players []string, startCount int, moves []Move, expectedTurnSequence []int, expectedWinner string) {
		var board Board
		board.Players = players
		board.PaperclipCount = startCount
		Play := func(idx int, m *Move) {
			prevPlayer := board.NextPlayer
			if err := board.Apply(m); err != nil {
				t.Error("Failed to apply move", idx, ":", m, ":", err.Error())
			}
			if board.NextPlayer == prevPlayer {
				t.Error("Failed to advance player counter?!")
			}
		}

		TestTurnSequence := func(expected, actual []int) {
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

		turnSequence := []int{0}
		for i, m := range moves {
			Play(i, &m)
			turnSequence = append(turnSequence, board.NextPlayer)
		}

		TestTurnSequence(expectedTurnSequence, turnSequence)

		if expectedWinner != "" && expectedWinner != board.Winner() {
			t.Error("Expected", expectedWinner, "to win, but", board.Winner(), "won")
		}
	}

	TestMoveSequence([]string{"a", "b"}, 5,
		[]Move{1, 2, 1, 1}, []int{0, 1, 0, 1, 0}, "b")

	TestMoveSequence([]string{"a", "b", "c"}, 5,
		[]Move{1, 2, 1, 1}, []int{0, 1, 2, 0, 1}, "a")

	TestMoveSequence([]string{"b", "a"}, 4,
		[]Move{1, 1, 2}, []int{0, 1, 0, 1}, "b")

	TestMoveSequence([]string{"b", "a"}, 4,
		[]Move{1, 1, 1, 1}, []int{0, 1, 0, 1, 0}, "a")

}
