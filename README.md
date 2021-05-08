# Golang finite state machine for decoding data packets

Inspired by Ragel, this is an attempt to create a way to easily define a data packets
for decoding incoming data with the ability to dynmicaly change the data definations
unlike the excellent, feature Ragel which generates source code from a data defiantion.

Example usage:

```
	// The machine will set the text variable
	var text string

	// Each state has a unique name, and each state has one or more transistions
	const (
		Start = "Start"
		Text  = "Text"
	)

	machine := Machine{
		// You must define the initial state of the machine
		InitialState: Start,
		
		// Here we define the states, and their transistions
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

	// Now we run the machine on some raw data
	n, err := machine.Parse([]byte{0x2, 'H', 'e', 'l', 'l', 'o', 0x0})
	// n will be set to the number of bytes used, and if err == nil then the data
	// was successfully parsed
```