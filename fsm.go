package fsm

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
)

type Machine struct {
	initialState      string
	States            map[string][]Transition
	logger            *log.Logger
	onErrorSkipBytes  int
	infiniteLoopCount int
}

type Transition struct {
	Test  TransitionTest
	State string
}

type TransitionTest func([]byte) (int, error)

type Option func(m *Machine)

func OptOnErrorSkipByte(v int) Option {
	return func(m *Machine) {
		m.onErrorSkipBytes = v
	}
}

func OptInfiniteLoopCount(v int) Option {
	return func(m *Machine) {
		m.infiniteLoopCount = v
	}
}

func OptLogger(logger *log.Logger) Option {
	return func(m *Machine) {
		m.logger = logger
	}
}

var ErrNoInitalState = errors.New("no inital state defined")
var ErrInfiniteLoop = errors.New("infinite loop detected")
var ErrNoTransitions = errors.New("no transitions")

func New(initalState string, states map[string][]Transition, options ...Option) *Machine {
	m := &Machine{
		initialState:      initalState,
		States:            states,
		onErrorSkipBytes:  0,
		infiniteLoopCount: 5000,
	}
	for _, opt := range options {
		opt(m)
	}
	return m
}

func (m *Machine) Parse(b []byte) (pos int, err error) {
	if m.logger == nil {
		m.logger = log.New(ioutil.Discard, "", 0)
	}
	defer func() {
		if err != nil {
			m.logger.Println(err.Error())
		}
	}()

	if m.infiniteLoopCount == 0 {
		m.infiniteLoopCount = 50000
	}

	state := m.initialState
	s, ok := m.States[state]
	if !ok {
		err = ErrNoInitalState
		return
	}

	counter := 0

RunState:
	counter++
	if m.onErrorSkipBytes == 0 && counter > m.infiniteLoopCount {
		return 0, fmt.Errorf("'%s': %w", state, ErrInfiniteLoop)
	}

	m.logger.Printf("entered state '%s' at position %d", state, pos)

	if len(s) == 0 {
		return 0, fmt.Errorf("'%s': %w", state, ErrNoTransitions)
	}

	var n int
	// Try each transition test, and if fails move onto next
	for i, t := range s {
		n, err = t.Test(b[pos:])
		if err != nil {
			m.logger.Printf(" - transition %d error: %s", i, err)
			continue
		}
		m.logger.Printf(" - transition %d used %d bytes", i, n)
		pos += n
		if t.State == "" {
			m.logger.Printf(" - SUCCESS: used %d bytes", pos)
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
