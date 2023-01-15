package jalex

import "fmt"

// ItemType represents the type of lexed item.
type ItemType string

// Common items.
const (
	ItemUnknown ItemType = "unknown"
	ItemError   ItemType = "error"
	ItemEOF     ItemType = "EOF"
)

// Item represents the structure sent to the client(e.g Parser).
type Item struct {
	T ItemType
	V string
}

func (i *Item) String() string {
	switch i.T {
	case ItemEOF:
		return string(ItemEOF)
	case ItemError:
		return i.V
	default:
		return fmt.Sprintf("%q", i.V)
	}
}
