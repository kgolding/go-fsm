package fsm

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
)

type Machine struct {
	InitialState string
	States       map[string][]Transition
	Logger       *log.Logger
}

type Transition struct {
	Test  TransitionTest
	State string
}

type TransitionTest func([]byte) (int, error)

type Error struct {
	State           string
	TransitionIndex int
	Err             error
}

func (e *Error) Error() string {
	return e.Err.Error()
}

var ErrNoInitalState = errors.New("no inital state defined")
var ErrInfiniteLoop = errors.New("infinite loop detected")

func (m *Machine) Parse(b []byte) (pos int, err error) {
	if m.Logger == nil {
		m.Logger = log.New(ioutil.Discard, "", 0)
	}
	defer func() {
		if err != nil {
			m.Logger.Println(err.Error())
		}
	}()

	state := m.InitialState
	s, ok := m.States[state]
	if !ok {
		err = ErrNoInitalState
		return
	}

	counter := 0

RunState:
	counter++
	if counter > 50000 {
		err = &Error{
			State:           state,
			TransitionIndex: -1,
			Err:             ErrInfiniteLoop,
		}
		return 0, ErrInfiniteLoop
	}

	m.Logger.Printf("entered state '%s' at position %d", state, pos)

	if len(s) == 0 {
		err = &Error{
			State:           state,
			TransitionIndex: -1,
			Err:             fmt.Errorf("state '%s' has no transitions!", state),
		}
		return
	}

	var n int
	// Try each transition test, and if fails move onto next
	for i, t := range s {
		n, err = t.Test(b[pos:])
		if err != nil {
			m.Logger.Printf(" - transition %d error: %s", i, err)
			continue
		}
		m.Logger.Printf(" - transition %d used %d bytes", i, n)
		pos += n
		if t.State == "" {
			m.Logger.Printf(" - SUCCESS: used %d bytes", pos)
			return // Success
		}
		// fmt.Println("Transition ", n, len(b))
		state = t.State
		var ok bool
		if s, ok = m.States[state]; ok {
			goto RunState
		}
		return 0, fmt.Errorf("no such state '%s'", state)
	}

	return // 0, fmt.Errorf("state '%s': no matching transitions", state)
}
