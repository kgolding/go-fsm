package fsm

import (
	"bytes"
	"fmt"
	"io"
)

// Skip ignores 1 or more bytes
func Skip(n int) TransitionTest {
	return func(b []byte) (int, error) {
		if len(b) < n {
			return 0, io.EOF
		}
		return n, nil
	}
}

// Else always changes to the given state
func Else() TransitionTest {
	return func(b []byte) (int, error) {
		return 0, nil
	}
}

// Byte matches a single given byte
func Byte(match byte) TransitionTest {
	return func(b []byte) (int, error) {
		if len(b) == 0 {
			return 0, io.EOF
		}
		if b[0] == match {
			return 1, nil
		}
		return 0, fmt.Errorf("looking for 0x%X got 0x%X", match, b[0])
	}
}

// Bytes matches a given []byte
func Bytes(match []byte) TransitionTest {
	return func(b []byte) (int, error) {
		if len(b) < len(match) {
			return 0, io.EOF
		}
		if bytes.Compare(b[:len(match)], match) == 0 {
			return len(match), nil
		}
		return 0, fmt.Errorf("looking for % X got % X", match, b[:len(match)])
	}
}

func STX() TransitionTest {
	return Byte(0x02)
}

func ETX() TransitionTest {
	return Byte(0x03)
}
