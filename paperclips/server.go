package paperclips

import (
	"errors"
	"fmt"
	"strconv"
)

type Mailbox struct {
	incoming chan BoardID
	poll     chan BoardID
	closer   chan bool
}

func NewMailbox() *Mailbox {
	ret := Mailbox{make(chan BoardID), make(chan BoardID), make(chan bool)}
	go ret.run()
	return &ret
}

func (b *Mailbox) run() {
	queued := make(map[BoardID]bool)
	defer close(b.incoming)
	defer close(b.poll)
	defer close(b.closer)

	for {
		// Depending on whether we have data to write out, use one of two selects
		//
		// TODO: There has to be a way to simplify this...
		add := func(id BoardID) {
			queued[id] = true
		}
		rm := func(id BoardID) {
			delete(queued, id)
		}

		if len(queued) > 0 {
			// If we have data, try to shove it out on the poll channel
			// while waiting for possible other data.
			var nextOut BoardID
			for nextOut, _ = range queued {
				break
			}

			select {
			case in := <-b.incoming:
				add(in)
			case b.poll <- nextOut:
				rm(nextOut)
			case <-b.closer:
				break
			}

		} else {
			// In the absence of ready outbound data, wait only on the input
			// side.
			select {
			case in := <-b.incoming:
				add(in)
			case <-b.closer:
				break
			}
		}
	}
}

func (b *Mailbox) Send(Id BoardID) {
	// Should return immediately, since the receiver is
	// busylooping. That is, if you set the thing up properly.
	b.incoming <- Id
}

func (b *Mailbox) Poll() (BoardID, bool) {
	// Wait until some data comes out. Or until the channel closes.
	data, ok := <-b.poll
	return data, ok
}

func (b *Mailbox) Shutdown() {
	// Should return immediately, since the receiver is
	// busylooping. That is, if you set the thing up properly.
	b.closer <- true
}

// TODO: make interface?
type Server struct {
	games       map[PlayerID]map[BoardID]*Board
	nextBoardId uint64
	mailboxes   map[PlayerID]*Mailbox
}

func NewServer() *Server {
	return &Server{make(map[PlayerID]map[BoardID]*Board), 0,
		make(map[PlayerID]*Mailbox)}
}

func (s *Server) Shutdown() {
	// TODO: Defer calls somehow instead?
	for _, box := range s.mailboxes {
		box.Shutdown()
	}
}

func (s *Server) PlayerExists(P PlayerID) bool {
	_, ret := s.games[P]
	return ret
}

func (s *Server) NewPlayer(Name PlayerID) error {
	if s.PlayerExists(Name) {
		return errors.New("Player " + string(Name) + " already exists on server")
	}
	s.games[Name] = make(map[BoardID]*Board)
	return nil
}

func (s *Server) GetPlayerList() []PlayerID {
	ret := make([]PlayerID, len(s.games))
	for p := range s.games {
		ret = append(ret, p)
	}
	return ret
}

func (s *Server) getNextBoardId() BoardID {
	ret := BoardID(strconv.FormatUint(s.nextBoardId, 10))
	s.nextBoardId++
	return ret
}

func (s *Server) NewGame(Players []PlayerID, StartCount int) (BoardID, error) {
	for _, p := range Players {
		if !s.PlayerExists(p) {
			//return errors.New("Player " + string(p) + " does not exist on server")
			if err := s.NewPlayer(p); err != nil {
				return "", err
			}
		}
	}

	board := NewBoard(Players, StartCount, s.getNextBoardId())
	ID := board.ID
	for _, p := range Players {
		s.games[p][ID] = board
	}
	return ID, nil
}

func (s *Server) GetGames(P PlayerID) (map[BoardID]*Board, error) {
	if !s.PlayerExists(P) {
		return nil, errors.New(fmt.Sprint("Player", P, "does not exist."))
	}
	return s.games[P], nil
}

func (s *Server) MakeMove(player PlayerID, board BoardID, move Move) error {
	if !s.PlayerExists(player) {
		return errors.New("Invalid player")
	}

	targetBoard, boardExists := s.games[player][board]
	if !boardExists {
		return errors.New("Invalid board")
	}

	if valid, err := Valid(&move, targetBoard); !valid {
		return errors.New("Invalid move: " + err.Error())
	}

	return nil
}
