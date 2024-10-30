package keymap

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

// Bag is a collection of keymaps
type Bag struct {
	keymaps []help.KeyMap
}

var _ help.KeyMap = Bag{}

// NewBag creates a new keymap bag
func NewBag(keymaps ...help.KeyMap) Bag {
	return Bag{keymaps: keymaps}
}

// Add adds keymaps to the bag
func (k Bag) Add(keymaps ...help.KeyMap) Bag {
	k.keymaps = append(k.keymaps, keymaps...)
	return k
}

// ShortHelp returns a slice of bindings to be displayed in the short version of the help
func (k Bag) ShortHelp() []key.Binding {
	var bindings []key.Binding
	for _, keymap := range k.keymaps {
		bindings = append(bindings, keymap.ShortHelp()...)
	}
	return bindings
}

// FullHelp returns an extended group of help items, grouped by columns
func (k Bag) FullHelp() [][]key.Binding {
	var bindings [][]key.Binding
	for _, keymap := range k.keymaps {
		bindings = append(bindings, keymap.FullHelp()...)
	}
	return bindings
}
