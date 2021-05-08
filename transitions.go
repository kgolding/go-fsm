package fsm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"regexp"
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

// Byte matches a single given byte
func Byte(match byte) TransitionTest {
	return func(b []byte) (int, error) {
		if len(b) == 0 {
			return 0, io.EOF
		}
		if b[0] == match {
			return 1, nil
		}
		return 0, errors.New("no null terminator")
	}
}

// RegexSubmatch matches a given regexp using FindSubmatch()
func RegexSubmatch(reg *regexp.Regexp, result *[][]byte) TransitionTest {
	return func(b []byte) (int, error) {
		r := reg.FindSubmatch(b)
		if r != nil {
			*result = r
			return len(r[0]), nil
		}
		return 0, errors.New("no match")
	}
}

// StringNullTerminated return a string of all bytes up to the first 0x00 / null byte
func StringNullTerminated(s *string) TransitionTest {
	return func(b []byte) (int, error) {
		p := bytes.IndexByte(b, 0x00)
		if p > -1 {
			*s = string(b[:p])
			return p + 1, nil
		}
		return 0, errors.New("No null in data")
	}
}

func Int16(v *int16) TransitionTest {
	return func(b []byte) (int, error) {
		if len(b) < 2 {
			return 0, io.EOF
		}
		*v = int16(binary.BigEndian.Uint16(b))
		return 2, nil
	}
}
