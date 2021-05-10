package fsm

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

type Scanner struct {
	machine     *Machine
	reader      io.Reader
	buf         []byte
	state       string
	transitions []Transition
	Err         error
}

func (m *Machine) NewScanner(reader io.Reader) *Scanner {
	s := &Scanner{
		machine:     m,
		reader:      reader,
		buf:         make([]byte, 0, 256),
		state:       m.initialState,
		transitions: m.States[m.initialState],
		Err:         nil,
	}
	return s
}

// Next blocks and waits for a success or the reader closing
func (s *Scanner) Next() bool {
	if s.machine.logger == nil {
		s.machine.logger = log.New(ioutil.Discard, "", 0)
	}

	if s.machine.infiniteLoopCount == 0 {
		s.machine.infiniteLoopCount = 50000
	}

	s.machine.logger.Println("Next()")

	b := make([]byte, 128)
	state := s.state
	counter := 0
	for {
		if len(s.buf) > 0 {
		RunState:
			counter++
			if s.machine.onErrorSkipBytes == 0 && counter > s.machine.infiniteLoopCount {
				s.Err = fmt.Errorf("'%s': %w", state, ErrInfiniteLoop)
				return false
			}

			s.machine.logger.Printf(" - trying state '%s'", state)

			if len(s.transitions) == 0 {
				s.Err = fmt.Errorf("'%s': %w", state, ErrNoTransitions)
				return false
			}

			softErrorFlag := false // If a transistion fails due to not enough bytes
			// Try each transition test, and if fails move onto next
			for i, t := range s.transitions {
				n, err := t.Test(s.buf)
				if err != nil {
					if err == io.EOF {
						softErrorFlag = true
					}
					s.machine.logger.Printf("   - transition %d error: %s", i, err)
					continue // next transition
				}
				s.machine.logger.Printf("   - transition %d used %d bytes [% X]", i, n, s.buf[:n])
				s.buf = s.buf[n:]
				if t.State == "" {
					s.machine.logger.Printf(" - SUCCESS: used %d bytes", n)
					s.state = s.machine.initialState
					s.transitions = s.machine.States[s.state]
					return true // Success
				}
				state = t.State
				var ok bool
				if s.transitions, ok = s.machine.States[state]; ok {
					goto RunState
				}
				s.Err = fmt.Errorf("no such state '%s'", state)
				return false
			}
			if !softErrorFlag {
				if s.machine.onErrorSkipBytes > 0 && len(s.buf) >= s.machine.onErrorSkipBytes {
					s.machine.logger.Printf("   - hard error, skipping %d byte(s)", s.machine.onErrorSkipBytes)
					s.buf = s.buf[s.machine.onErrorSkipBytes:]
				}
			}
		}
		n, err := s.reader.Read(b)
		if errors.Is(err, io.EOF) && len(s.buf) > 0 {
			// The reader has closed but we still have data in the buffer to process
			continue
		}
		if err != nil {
			s.machine.logger.Printf(" - read error: %s (p.buf still has %d bytes)", err, len(s.buf))
			s.Err = err
			return false
		}
		s.buf = append(s.buf, b[:n]...)
	}
}
