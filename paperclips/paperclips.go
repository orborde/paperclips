package paperclips

import (
	"errors"
	"fmt"
	"strconv"
)

type PaperclipCount int
type Move PaperclipCount

const (
	TakeOne Move = 1
	TakeTwo Move = 2
)

type PlayerID string
type BoardID string
type TurnCount int

type Board struct {
	PaperclipCount PaperclipCount
}

type MoveResult struct {
	*BoardMessage
	error
}

type MoveMessage struct {
	Move Move
	// TODO: Do we need this for verification?
	// TODO: How do we communicate errors? I suppose we can just crash out for now...
	Player PlayerID
	// Turn at which this Move was issued. It might get dropped on the
	// floor if it's too old.
	TurnCount TurnCount
	// Channel to signal back the move results
	Result chan MoveResult
}

type BoardMessage struct {
	Board     Board
	WhoseTurn PlayerID
	TurnCount TurnCount
}

func Play(Players []PlayerID, StartCounter PaperclipCount,
	Moves <-chan MoveMessage, End <-chan bool, Updates chan<- BoardMessage) {
	// Set up the initial game state.
	currentPlayerIndex := 0
	board := NewBoard(StartCounter)
	// TODO: redundant with board counter
	var turnCount TurnCount = 0

	currentStatus := func() *BoardMessage {
		return &BoardMessage{*board, Players[currentPlayerIndex], turnCount}
	}

	Updates <- *currentStatus()

	applyMove := func(move *MoveMessage) {
		if move.TurnCount < turnCount {
			move.Result <- MoveResult{
				nil, errors.New("Move discarded because it refers to an out of date board.")}
			return
		}

		// TODO: check for correct player

		err := board.Apply(&move.Move)
		if err != nil {
			move.Result <- MoveResult{nil, errors.New(fmt.Sprint("Move rejected as invalid: ", err))}
			return
		}

		turnCount++
		currentPlayerIndex++
		currentPlayerIndex %= len(Players)

		Updates <- *currentStatus()
		move.Result <- MoveResult{currentStatus(), nil}
	}

	select {
	case <-End:
		return
	case move := <-Moves:
		applyMove(&move)
	}
}

func NewBoard(StartCount PaperclipCount) *Board {
	return &Board{PaperclipCount: StartCount}
}

func (b *Board) Render() string {
	return strconv.Itoa(int(b.PaperclipCount))
}

func (b *Board) Apply(move *Move) error {
	if valid, err := Valid(move, b); !valid {
		return errors.New(fmt.Sprint("Could not apply move ", move, ":", err.Error()))
	}
	b.PaperclipCount -= PaperclipCount(*move)

	return nil
}

func (b *Board) GameOver() bool {
	return b.PaperclipCount == 0
}

func (m *Move) Valid() (bool, error) {
	if *m != TakeOne && *m != TakeTwo {
		return false, errors.New("Not a valid move type")
	}
	return true, nil
}

func Valid(m *Move, b *Board) (bool, error) {
	if valid, err := m.Valid(); !valid {
		return false, err
	}
	if b.GameOver() {
		return false, errors.New("Game has already ended")
	}
	if PaperclipCount(*m) > b.PaperclipCount {
		return false, errors.New("Not enough paperclips")
	}
	return true, nil
}
