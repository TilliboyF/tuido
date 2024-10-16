package tui

import "github.com/charmbracelet/bubbles/key"


type Keys struct {
	LineUp   key.Binding
	LineDown key.Binding
	Edit     key.Binding
	New      key.Binding
	Delete   key.Binding
	View     key.Binding
	Quit     key.Binding
}

/* Implementation for the Keymap interface */
func (km Keys) ShortHelp() []key.Binding {
	return []key.Binding{km.LineUp, km.LineDown, km.Edit, km.New, km.Delete, km.View, km.Quit}
}

// FullHelp implements the KeyMap interface.
func (km Keys) FullHelp() [][]key.Binding {
	return nil
}

func defaultKeys() Keys {
	return Keys{
		LineUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		LineDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
		New: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
		View: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "view todo"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
	}
}
