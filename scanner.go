package fsm

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

type Scanner struct {
	machine *Machine
	reader  io.Reader
	buf     []byte
	state   string
	s       []Transition
	Err     error
}

func (m *Machine) NewScanner(reader io.Reader) *Scanner {
	return &Scanner{
		machine: m,
		reader:  reader,
		buf:     make([]byte, 0, 256),
		state:   m.InitialState,
		s:       m.States[m.InitialState],
		Err:     nil,
	}
}

// Next blocks and waits for a success or the reader closing
func (p *Scanner) Next() bool {
	if p.machine.Logger == nil {
		p.machine.Logger = log.New(ioutil.Discard, "", 0)
	}

	b := make([]byte, 128)
	state := p.state
	for {
		if len(p.buf) > 0 {
			counter := 0
		RunState:
			counter++
			if counter > 50000 {
				p.Err = &Error{
					State:           state,
					TransitionIndex: -1,
					Err:             ErrInfiniteLoop,
				}
				return false
			}

			p.machine.Logger.Printf("trying state '%s'", state)

			if len(p.s) == 0 {
				p.Err = &Error{
					State:           state,
					TransitionIndex: -1,
					Err:             fmt.Errorf("state '%s' has no transitions!", state),
				}
				return false
			}

			// Try each transition test, and if fails move onto next
			for i, t := range p.s {
				n, err := t.Test(p.buf)
				if err != nil {
					p.machine.Logger.Printf(" - transition %d error: %s", i, err)
					continue // next transition
				}
				p.machine.Logger.Printf(" - transition %d used %d bytes", i, n)
				p.buf = p.buf[n:]
				if t.State == "" {
					p.machine.Logger.Printf(" - SUCCESS: used %d bytes", n)
					p.state = p.machine.InitialState
					p.s = p.machine.States[p.state]
					return true // Success
				}
				// fmt.Println("Transition ", n, len(b))
				state = t.State
				var ok bool
				if p.s, ok = p.machine.States[state]; ok {
					goto RunState
				}
				p.Err = fmt.Errorf("no such state '%s'", state)
				return false
			}
		}
		n, err := p.reader.Read(b)
		if err != nil {
			p.Err = err
			return false
		}
		p.buf = append(p.buf, b[:n]...)
	}
}
