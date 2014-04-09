package paperclips

import "testing"

func TestRender(t *testing.T) {
	const render = ("X|O| " +
		"-----" +
		" | | " +
		"-----" +
		"X|X|O")
	data := Board{
		{X, O, Blank},
		{Blank, Blank, Blank},
		{X, X, O},
	}

	if data.Render() != render {
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

	ExpectValid(Move{0, 0, X})
	ExpectValid(Move{0, 0, O})
	ExpectInvalid(Move{0, 0, Blank})
	ExpectInvalid(Move{-1, 0, X})
	ExpectInvalid(Move{0, -1, X})
	ExpectInvalid(Move{-1, -1, X})
	ExpectInvalid(Move{0, 3, X})
	ExpectInvalid(Move{3, 0, X})
}
