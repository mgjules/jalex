package jalex_test

import (
	"strings"
	"testing"

	"github.com/mgjules/jalex"
	"github.com/stretchr/testify/assert"
)

const (
	itemText      jalex.ItemType = "text"
	itemLeftMeta  jalex.ItemType = "left_meta"
	itemRightMeta jalex.ItemType = "right_meta"
	itemUser      jalex.ItemType = "user"
)

func TestLexer(t *testing.T) {
	l, err := jalex.NewLexer("TestLexer", textStateFn)
	assert.NoError(t, err)
	assert.NotNil(t, l)

	var items []jalex.Item
	itemsCh := l.Run("Hello <@mike>! How are you?")
	for item := range itemsCh {
		items = append(items, item)
	}

	assert.Equal(t, []jalex.Item{
		{itemText, "Hello "},
		{itemLeftMeta, "<"},
		{itemUser, "mike"},
		{itemRightMeta, ">"},
		{itemText, "! How are you?"},
		{jalex.ItemEOF, ""},
	}, items)
}

func textStateFn(s *jalex.State) jalex.StateFn {
	for {
		if strings.HasPrefix(s.Text[s.Current:], "<") {
			if s.Current > s.Start {
				s.Emit(itemText)
			}
			return leftMetaStateFn
		}
		if s.Next() == jalex.EOF {
			break
		}
	}
	if s.Current > s.Start {
		s.Emit(itemText)
	}
	s.Emit(jalex.ItemEOF)
	return nil
}

func leftMetaStateFn(s *jalex.State) jalex.StateFn {
	s.Current += len("<")
	s.Emit(itemLeftMeta)
	return insideMetaStateFn
}

func insideMetaStateFn(s *jalex.State) jalex.StateFn {
	for {
		if strings.HasPrefix(s.Text[s.Current:], ">") {
			return rightMetaStateFn
		}

		r := s.Next()
		switch {
		case r == jalex.EOF || r == '\n':
			return s.Errorf("unclosed meta")
		case r == '@':
			s.Skip() //ignore the '@' character
			return userStateFn
		}
	}
}

func rightMetaStateFn(s *jalex.State) jalex.StateFn {
	s.Current += len(">")
	s.Emit(itemRightMeta)
	return textStateFn
}

func userStateFn(s *jalex.State) jalex.StateFn {
	for strings.IndexRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", s.Next()) >= 0 {
	}
	s.Back()
	s.Emit(itemUser)
	return insideMetaStateFn
}
