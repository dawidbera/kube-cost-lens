// Package tui provides keyboard shortcuts and help definitions.
package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the full set of interactive keybindings.
type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	NextTab  key.Binding
	PrevTab  key.Binding
	Enter    key.Binding
	Patch    key.Binding
	Filter   key.Binding
	Escape   key.Binding
	Copy     key.Binding
	Write    key.Binding
	Refresh  key.Binding
	Help     key.Binding
	Quit     key.Binding
}

// DefaultKeyMap returns standard keyboard mapping for the application.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		NextTab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next view"),
		),
		PrevTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev view"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "inspect containers"),
		),
		Patch: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "view patch"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter search"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "close / clear"),
		),
		Copy: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "copy patch"),
		),
		Write: key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "save patch file"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}
