package fsm

import (
	"errors"
	"regexp"
)

var ErrRegexNoMatch = errors.New("regex no match")

// RegexFind matches a given regexp using Find()
func RegexFind(reg *regexp.Regexp, result *[]byte) TransitionTest {
	return func(b []byte) (int, error) {
		r := reg.Find(b)
		if r != nil {
			*result = r
			return len(r), nil
		}
		return 0, ErrRegexNoMatch
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
		return 0, ErrRegexNoMatch
	}
}
