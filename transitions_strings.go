package fsm

import (
	"bytes"
	"fmt"
	"io"
)

func StringLen(length *uint, s *string) TransitionTest {
	return func(b []byte) (int, error) {
		l := int(*length)
		if len(b) < l {
			return 0, io.EOF
		}

		*s = string(b[:l])
		return l, nil
	}
}

func StringDelimited(delimiter byte, s *string) TransitionTest {
	return func(b []byte) (int, error) {
		p := bytes.IndexByte(b, delimiter)
		if p > -1 {
			*s = string(b[:p])
			return p + 1, nil
		}
		return 0, fmt.Errorf("no delimted %X in data", delimiter)
	}
}

// StringNullTerminated return a string of all bytes up to the first 0x00 / null byte
func StringNullTerminated(s *string) TransitionTest {
	return StringDelimited(0x0, s)
}