package jalex

import (
	"fmt"
	"unicode/utf8"
)

const EOF = -1

// State holds the current state of the currently running state machine.
type State struct {
	Text    string
	Start   int // the beginning cursor position of a selection of the Text field
	Current int // the current cursor position of a selection of the Text field
	Width   int
	items   chan Item
}

// Emit send an item on the items channel and move s.start to s.current.
func (s *State) Emit(t ItemType) {
	s.items <- Item{t, s.Text[s.Start:s.Current]}
	s.Start = s.Current
}

// Errorf send an ItemError on the items channel and returns an nil StateFn.
func (s *State) Errorf(format string, args ...interface{}) StateFn {
	format = fmt.Sprintf("[Current Position: %d] %s", s.Current, format)
	s.items <- Item{ItemError, fmt.Sprintf(format, args...)}
	return nil
}

// Next moves the current cursor position forward by one rune.
func (s *State) Next() rune {
	if s.Current >= len(s.Text) {
		s.Width = 0
		return EOF
	}

	var r rune
	r, s.Width = utf8.DecodeRuneInString(s.Text[s.Current:])
	s.Current += s.Width

	return r
}

// Skip skips the current selection of text.
func (s *State) Skip() {
	s.Start = s.Current
}

// Back moves the current cursor position backward by one rune.
func (s *State) Back() {
	s.Current -= s.Width
}

// StateFn function represents one state function.
type StateFn func(*State) StateFn
