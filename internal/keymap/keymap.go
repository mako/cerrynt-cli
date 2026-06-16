// Package keymap defines the default keybindings used across all screens.
// The bubbles/key package integrates with Bubble Tea's help system and
// provides a clean way to match key presses in Update() functions.
package keymap

import "github.com/charmbracelet/bubbles/key"

// KeyMap holds all keybindings for the application.
type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Open   key.Binding
	Back   key.Binding
	Quit   key.Binding
}

// Default is the application-wide keymap with vim-style defaults.
var Default = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j", "down"),
	),
	Open: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "open"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
