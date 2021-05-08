package fsm

import (
	"log"
	"os"
	"regexp"
	"testing"
)

var logger = log.New(os.Stdout, "FSM ", log.Lmsgprefix)

func init() {
	logger = nil // Comment this out to enable logging
}

func TestDecodeSimple(t *testing.T) {
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

	machine.Logger = logger

	n, err := machine.Parse([]byte{0x2, 'H', 'e', 'l', 'l', 'o', 0x0})
	if err != nil {
		t.Error(err)
	}
	if text != "Hello" {
		t.Errorf("Expected ' Hello', got '%s'", text)
	}
	if n != 7 {
		t.Errorf("Expected 7, got %d", n)
	}
}

func TestDecode1(t *testing.T) {
	var text string
	var dateMatch [][]byte

	const (
		Start = "Start"
		Date  = "Date"
		Text  = "Text"
	)

	machine := Machine{
		InitialState: Start,
		States: map[string][]Transition{
			Start: []Transition{
				{Byte(0x0A), Date},
				{Skip(1), Start},
			},
			Date: []Transition{
				{RegexSubmatch(regexp.MustCompile("([0-9]{2})/([0-9]{2})/([0-9]{2,4})"), &dateMatch), Text},
			},
			Text: []Transition{
				{StringNullTerminated(&text), ""},
			},
		},
	}

	machine.Logger = logger

	n, err := machine.Parse(append([]byte("X\n08/05/2021 Hello!"), 0x00))
	if err != nil {
		t.Error(err)
	}
	if len(dateMatch) != 4 {
		t.Errorf("Expected 4, got %d", len(dateMatch))
	} else {
		if string(dateMatch[0]) != "08/05/2021" {
			t.Errorf("Expected '08/05/2021', got '%s'", dateMatch[0])
		}
	}
	if text != " Hello!" {
		t.Errorf("Expected ' Hello!', got '%s'", text)
	}
	if n != 20 {
		t.Errorf("Expected 20, got %d", n)
	}
}

func TestDecodeVariableLenString(t *testing.T) {
	var text string
	var textLen uint

	const (
		Start   = "Start"
		TextLen = "TextLen"
		Text    = "Text"
	)

	machine := Machine{
		InitialState: Start,
		States: map[string][]Transition{
			Start: []Transition{
				{STX(), TextLen},
				{Skip(1), Start},
			},
			TextLen: []Transition{
				{Uint(&textLen, 2), Text},
			},
			Text: []Transition{
				{StringLen(&textLen, &text), ""},
			},
		},
	}

	machine.Logger = logger

	n, err := machine.Parse([]byte{0x2, 0x0, 0x5, 'H', 'e', 'l', 'l', 'o'})
	if err != nil {
		t.Error(err)
	}
	if text != "Hello" {
		t.Errorf("Expected ' Hello', got '%s'", text)
	}
	if n != 8 {
		t.Errorf("Expected 8, got %d", n)
	}
}