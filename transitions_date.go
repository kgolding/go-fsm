package fsm

import (
	"io"
	"time"
)

// Date reads a date string using the time.Parse layout Mon Jan 2 15:04:05 -0700 MST 2006
func StringDate(layout string, t *time.Time) TransitionTest {
	return func(b []byte) (int, error) {
		if len(b) < len(layout) {
			return 0, io.EOF
		}
		var err error
		*t, err = time.Parse(layout, string(b))
		if err != nil {
			return 0, err
		}
		return len(layout), nil
	}
}
