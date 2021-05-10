package fsm

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

type Scanner struct {
	machine          *Machine
	reader           io.Reader
	onErrorSkipBytes int
	buf              []byte
	state            string
	s                []Transition
	Err              error
}

type Option func(s *Scanner)

func OptOnErrorSkipByte(v int) Option {
	return func(s *Scanner) {
		s.onErrorSkipBytes = v
	}
}

func (m *Machine) NewScanner(reader io.Reader, options ...Option) *Scanner {
	s := &Scanner{
		machine: m,
		reader:  reader,
		buf:     make([]byte, 0, 256),
		state:   m.InitialState,
		s:       m.States[m.InitialState],
		Err:     nil,
	}
	for _, opt := range options {
		opt(s)
	}
	return s
}

// Next blocks and waits for a success or the reader closing
func (p *Scanner) Next() bool {
	if p.machine.Logger == nil {
		p.machine.Logger = log.New(ioutil.Discard, "", 0)
	}

	p.machine.Logger.Println("Next()")

	b := make([]byte, 128)
	state := p.state
	counter := 0
	for {
		if len(p.buf) > 0 {
		RunState:
			counter++
			if counter > 500 {
				p.Err = fmt.Errorf("'%s': %w", state, ErrInfiniteLoop)
				return false
			}

			p.machine.Logger.Printf(" - trying state '%s'", state)

			if len(p.s) == 0 {
				p.Err = fmt.Errorf("'%s': %w", state, ErrNoTransitions)
				return false
			}

			softErrorFlag := false // If a transistion fails due to not enough bytes
			// Try each transition test, and if fails move onto next
			for i, t := range p.s {
				n, err := t.Test(p.buf)
				if err != nil {
					if err == io.EOF {
						softErrorFlag = true
					}
					p.machine.Logger.Printf("   - transition %d error: %s", i, err)
					continue // next transition
				}
				p.machine.Logger.Printf("   - transition %d used %d bytes [% X]", i, n, p.buf[:n])
				p.buf = p.buf[n:]
				if t.State == "" {
					p.machine.Logger.Printf(" - SUCCESS: used %d bytes", n)
					p.state = p.machine.InitialState
					p.s = p.machine.States[p.state]
					return true // Success
				}
				state = t.State
				var ok bool
				if p.s, ok = p.machine.States[state]; ok {
					goto RunState
				}
				p.Err = fmt.Errorf("no such state '%s'", state)
				return false
			}
			if !softErrorFlag {
				if p.onErrorSkipBytes > 0 && len(p.buf) >= p.onErrorSkipBytes {
					p.machine.Logger.Printf("   - hard error, skipping %d byte(s)", p.onErrorSkipBytes)
					p.buf = p.buf[p.onErrorSkipBytes:]
				}
			}
		}
		n, err := p.reader.Read(b)
		if errors.Is(err, io.EOF) && len(p.buf) > 0 {
			// The reader has closed but we still have data in the buffer to process
			continue
		}
		if err != nil {
			p.machine.Logger.Printf(" - read error: %s (p.buf still has %d bytes)", err, len(p.buf))
			p.Err = err
			return false
		}
		p.buf = append(p.buf, b[:n]...)
	}
}
