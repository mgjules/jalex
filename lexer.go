package jalex

import (
	"crypto/sha1"
	"errors"
	"fmt"

	"golang.org/x/sync/singleflight"
)

// Lexer holds the required information to run several state machines.
type Lexer struct {
	name           string
	initialStateFn StateFn            // initial state function
	antiDup        singleflight.Group // prevents duplicate "Run"
}

// NewLexer creates a new Lexer with the given name and an initial state function.
func NewLexer(name string, initial StateFn) (*Lexer, error) {
	if name == "" {
		return nil, errors.New("name must not be empty")
	}

	if initial == nil {
		return nil, errors.New("initial state function must not be nil")
	}

	return &Lexer{
		name:           name,
		initialStateFn: initial,
	}, nil
}

// Run runs the state machine for a given text string and returns a channel of Item.
func (l *Lexer) Run(text string) chan Item {
	state := &State{
		Text:  text,
		items: make(chan Item),
	}

	hashed := fmt.Sprintf("%x", sha1.Sum([]byte(text)))

	go l.antiDup.Do(hashed, func() (interface{}, error) {
		defer close(state.items)

		for stateFn := l.initialStateFn; stateFn != nil; {
			stateFn = stateFn(state)
		}

		return nil, nil
	})

	return state.items
}
