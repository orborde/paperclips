package tictactoe

import "testing"

func TestRender(t *testing.T) {
	const render = (
		"X|O| " +
		"-----" +
		" | | " +
		"-----" +
		"X|X|O")
	data := Map{
		{X, O, Blank},
		{Blank, Blank, Blank},
		{X, X, O},
	}

	if data.Render() != render {
		t.Errorf("Failed to properly render board!")
	}
}
