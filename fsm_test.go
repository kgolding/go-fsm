package fsm

import (
	"log"
	"os"
	"regexp"
	"testing"
)

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
		States: map[string]*State{
			Start: &State{
				Transitions: []Transition{
					{Byte(0x0A), Date},
					{Skip(1), Start},
				},
			},
			Date: &State{
				Transitions: []Transition{
					{RegexSubmatch(regexp.MustCompile("([0-9]{2})/([0-9]{2})/([0-9]{2,4})"), &dateMatch), Text},
				},
			},
			Text: &State{
				Transitions: []Transition{
					{StringNullTerminated(&text), ""},
				},
			},
		},
	}

	machine.Logger = log.New(os.Stdout, "FSM ", log.Lmsgprefix)

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
