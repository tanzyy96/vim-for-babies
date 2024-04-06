package main

import (
	"github.com/charmbracelet/bubbles/key"
)

// listKeymap is the "class" for keybindings for ListView
// Put the additional keymaps here. Default includes up, down, Home, End, Quit
type listKeymap struct {
	AddMotion key.Binding
}

func (km *listKeymap) Bindings() []key.Binding {
	return []key.Binding{km.AddMotion}
}

// WithKeys selects the keys being targeted by the binding. You can select multiple keys per binding.
// WithHelp provides the help text for given keybinding
func newListKeymap() additionalKeymap {
	return &listKeymap{
		AddMotion: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "Add a new motion"),
		),
	}
}

type addMotionKeymap struct {
	Up, Down, Submit, Esc key.Binding
}

func (km *addMotionKeymap) Bindings() []key.Binding {
	return []key.Binding{km.Up, km.Down, km.Submit, km.Esc}
}

func (km *addMotionKeymap) ShortHelp() []key.Binding {
	return km.Bindings()
}

func (km *addMotionKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{km.Up, km.Down},
		{km.Submit, km.Esc},
	}
}

func newAddMotionKeymap() additionalKeymap {
	return &addMotionKeymap{
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("up", "Up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("down", "Down"),
		),
		Submit: key.NewBinding(
			key.WithKeys("shift+enter"),
			key.WithHelp("shift+enter", "Submit"),
		),
		Esc: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "Back"),
		),
	}
}
