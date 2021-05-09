package fsm

import (
	"bytes"
	"io"
	"testing"
	"time"
)

func Test_Parser(t *testing.T) {
	var text string
	const (
		Start = "Start"
		Text  = "Text"
	)
	machine := Machine{
		InitialState: Start,
		States: map[string][]Transition{
			Start: []Transition{
				{STX(), Text},
				{Skip(1), Start},
			},
			Text: []Transition{
				{StringNullTerminated(&text), ""},
			},
		},
	}
	// machine.Logger = logger

	// Mock an io.Reader with slow data being received
	count := 50
	r, w := io.Pipe()
	hello := []byte{0x2, 'H', 'e', 'l', 'l', 'o', 0x0}
	go func() {
		b := bytes.Repeat(hello, count)
		for len(b) > 0 {
			l := 2
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
