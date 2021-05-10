package fsm

import (
	"bytes"
	"io"
	"testing"
	"time"
)

func Test_Scanner(t *testing.T) {
	var text string
	const (
		Start = "Start"
		Text  = "Text"
	)
	machine := New(
		Start,
		map[string][]Transition{
			Start: []Transition{
				{STX(), Text},
				{Skip(1), Start},
			},
			Text: []Transition{
				{StringNullTerminated(&text), ""},
			},
		},
		// OptLogger(logger),
	)

	// Mock an io.Reader with slow data being received, test sending 1, 2, 4, 8 ... 256 bytes at a time
	for batchLen := 1; batchLen < 256; batchLen *= 2 {
		count := 50
		r, w := io.Pipe()
		hello := []byte{0x2, 'H', 'e', 'l', 'l', 'o', 0x0}
		go func() {
			b := bytes.Repeat(hello, count)
			for len(b) > 0 {
				l := batchLen
				if len(b) < l {
					l = len(b)
				}
				w.Write(b[:l])
				b = b[l:]
				time.Sleep(time.Millisecond * 2)
			}
			w.Close()
		}()

		// p := machine.NewScanner(r, OptOnErrorSkipByte(1))
		p := machine.NewScanner(r)

		mcount := 0
		for p.Next() {
			if text != "Hello" {
				t.Fatalf("expected 'Hello' got '%s'", text)
			}
			mcount++
			// t.Log("MATCH", mcount, text)
			text = ""
		}
		if p.Err != io.EOF {
			t.Error(p.Err)
		}
		if mcount != count {
			t.Errorf("expected %d got %d matches", count, mcount)
		}
	}
}

func Test_ScannerLoop(t *testing.T) {
	var text string
	const (
		Start = "Start"
		Text  = "Text"
	)
	machine := New(
		Start,
		map[string][]Transition{
			Start: []Transition{
				{STX(), Text},
				{Skip(1), Start},
			},
			Text: []Transition{
				{StringNullTerminated(&text), ""},
			},
		},
		OptOnErrorSkipByte(1),
		// OptInfiniteLoopCount(500000),
		// OptLogger(logger),
	)

	// Mock an io.Reader with slow data being received, test sending 1, 2, 4, 8 ... 256 bytes at a time
	for batchLen := 1; batchLen < 256; batchLen *= 2 {
		count := 50
		r, w := io.Pipe()
		hello := []byte{0x2, 'H', 'e', 'l', 'l', 'o', 0x0}
		go func() {
			b := bytes.Repeat(hello, count-1)
			b = append(b, bytes.Repeat([]byte{0xff}, 5000)...)
			b = append(b, hello...)
			for len(b) > 0 {
				l := batchLen
				if len(b) < l {
					l = len(b)
				}
				w.Write(b[:l])
				b = b[l:]
			}
			w.Close()
		}()

		p := machine.NewScanner(r)

		mcount := 0
		for p.Next() {
			if text != "Hello" {
				t.Fatalf("expected 'Hello' got '%s'", text)
			}
			mcount++
			// t.Log("MATCH", mcount, text)
			text = ""
		}
		if p.Err != io.EOF {
			t.Error(p.Err)
		}
		if mcount != count {
			t.Errorf("expected %d got %d matches", count, mcount)
		}
	}
}

func Test_ScannerBadData(t *testing.T) {
	var text string
	var length uint
	var date time.Time

	const (
		Start = "Start"
		Len   = "Len"
		Text  = "Text"
		Date  = "Date"
	)
	machine := New(
		Start,
		map[string][]Transition{
			Start: []Transition{
				{STX(), Len},
			},
			Len: []Transition{
				{Uint(&length, 1), Text},
			},
			Text: []Transition{
				{StringFixedLen(&length, &text), Date},
			},
			Date: []Transition{
				{DateString("02/01/06", &date), ""},
			},
		},
		OptOnErrorSkipByte(1),
		// OptLogger(logger),
	)

	// Mock an io.Reader with slow data being received, test sending 1, 2, 4, 8 ... 256 bytes at a time
	for batchLen := 1; batchLen < 256; batchLen *= 2 {
		r, w := io.Pipe()
		hello := []byte{0x2, 0x05, 'H', 'e', 'l', 'l', 'o', '0', '1', '/', '1', '2', '/', '2', '1'} // 15 bytes
		b := bytes.Repeat(hello, 2)
		b = append(b, []byte{0xff, 0xff, 0xff}...) // Bad bytes between good data
		b = append(b, bytes.Repeat(hello, 2)...)
		// Bad date
		b = append(b, []byte{0x2, 0x05, 'H', 'e', 'l', 'l', 'o', '0', '1', '/', '1', '2', '/', 'X', '1'}...)
		b = append(b, []byte{0xff, 0xff, 0xff}...) // Bad bytes between good data
		b = append(b, bytes.Repeat(hello, 2)...)
		count := 6
		go func() {
			for len(b) > 0 {
				l := batchLen
				if len(b) < l {
					l = len(b)
				}
				w.Write(b[:l])
				b = b[l:]
				time.Sleep(time.Millisecond * 2)
			}
			w.Close()
		}()

		p := machine.NewScanner(r)

		mcount := 0
		for p.Next() {
			if text != "Hello" {
				t.Errorf("expected 'Hello' got '%s'", text)
			}
			if length != 5 {
				t.Errorf("expect 5 got %d", length)
			}
			mcount++
			// t.Log("MATCH", mcount, text)
			text = ""
			length = 0
		}
		if p.Err != io.EOF {
			t.Errorf("unexpected error: %s", p.Err)
		}
		if mcount != count {
			t.Errorf("expected %d got %d matches", count, mcount)
		}
	}
}
