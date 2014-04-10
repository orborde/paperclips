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
	Winner    *PlayerID
	TurnCount TurnCount
}

type Game struct {
	Moves       chan MoveMessage
	end         chan bool
	FirstUpdate chan BoardMessage
}

func NewGame(Players []PlayerID, StartCounter PaperclipCount) *Game {
	ret := &Game{make(chan MoveMessage), make(chan bool), make(chan BoardMessage)}
	go ret.Play(Players, StartCounter)
	return ret
}

func (g *Game) End() {
	g.end <- true
}

func (g *Game) Play(Players []PlayerID, StartCounter PaperclipCount) {
	// Set up the initial game state.
	currentPlayerIndex := 0
	board := NewBoard(StartCounter)
	// TODO: redundant with board counter
	var turnCount TurnCount = 0
	var winner *PlayerID = nil

	currentPlayer := func() *PlayerID {
		return &Players[currentPlayerIndex]
	}
	currentStatus := func() *BoardMessage {
		return &BoardMessage{*board, *currentPlayer(), winner, turnCount}
	}
	applyMove := func(move *MoveMessage) {
		if move.TurnCount < turnCount {
			move.Result <- MoveResult{
				nil, errors.New("Move discarded because it refers to an out of date board.")}
			return
		}

		// TODO: check for correct player
		// TODO: check for gameover invalidating move?

		err := board.Apply(&move.Move)
		if err != nil {
			move.Result <- MoveResult{nil, errors.New(fmt.Sprint("Move rejected as invalid: ", err))}
			return
		}

		turnCount++
		if board.GameOver() {
			winner = currentPlayer()
		}
		currentPlayerIndex++
		currentPlayerIndex %= len(Players)

		move.Result <- MoveResult{currentStatus(), nil}
	}

	g.FirstUpdate <- *currentStatus()
	close(g.FirstUpdate)

	defer close(g.end)
	defer close(g.Moves)
	for {
		select {
		case <-g.end:
			return
		case move := <-g.Moves:
			applyMove(&move)
		}
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
