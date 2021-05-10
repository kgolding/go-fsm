package fsm

import (
	"bytes"
	"fmt"
	"io"
)

// StringFixedLen reads a fixed length string
func StringFixedLen(length *uint, s *string) TransitionTest {
	return func(b []byte) (int, error) {
		l := int(*length)
		if len(b) < l {
			return 0, io.EOF
		}

		*s = string(b[:l])
		return l, nil
	}
}

// StringDelimited reads a string up to the given delimiter
func StringDelimited(delimiter byte, s *string) TransitionTest {
	return func(b []byte) (int, error) {
		p := bytes.IndexByte(b, delimiter)
		if p == -1 {
			return 0, io.EOF
		}
		*s = string(b[:p])
		return p + 1, nil
	}
}

// StringDelimited reads a string up to the given delimiter
func StringDelimitedMaxLen(delimiter byte, maxLen int, s *string) TransitionTest {
	return func(b []byte) (int, error) {
		p := bytes.IndexByte(b, delimiter)
		if p == -1 {
			if len(b) < maxLen {
				return 0, io.EOF
			}
			return 0, fmt.Errorf("no delimiter %X in data", delimiter)
		}
		*s = string(b[:p])
		return p + 1, nil
	}
}

// StringNullTerminated reads a string of all bytes up to the first 0x00 / null byte
func StringNullTerminated(s *string) TransitionTest {
	return StringDelimited(0x0, s)
}
