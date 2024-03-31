package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	items := []list.Item{
		item{title: "Move cursor", desc: "hjkl"},
		item{title: "Go to top of page", desc: "gg"},
		item{title: "Go to bottom of page", desc: "G"},
		item{title: "Move to next word", desc: "w/W"},
		item{title: "Move to previous word", desc: "b/B"},
		item{title: "Move to end of word", desc: "e/E"},
		item{title: "Move to beginning of line", desc: "0"},
		item{title: "Move to end of line", desc: "$"},
		item{title: "Move to first non-whitespace character of line", desc: "_"},
		item{title: "Move to top of screen", desc: "H"},
		item{title: "Move to middle of screen", desc: "M"},
		item{title: "Move to bottom of screen", desc: "L"},
		item{title: "Move up half a page", desc: "Ctrl+u"},
		item{title: "Move down half a page", desc: "Ctrl+d"},
		item{title: "Move up a page", desc: "Ctrl+b"},
		item{title: "Move down a page", desc: "Ctrl+f"},
		item{title: "Move screen up one line", desc: "Ctrl+e"},
		item{title: "Move screen down one line", desc: "Ctrl+y"},
		item{title: "Replace char", desc: "r"},
		item{title: "Replace line", desc: "R"},
		item{title: "Insert before cursor", desc: "i"},
		item{title: "Insert at beginning of line", desc: "I"},
		item{title: "Append after cursor", desc: "a"},
		item{title: "Append at end of line", desc: "A"},
		item{title: "Insert new line below cursor", desc: "o"},
		item{title: "Insert new line above cursor", desc: "O"},
		item{title: "Delete char", desc: "x"},
		item{title: "Delete line", desc: "dd"},
		item{title: "Delete word", desc: "dw"},
		item{title: "Delete to end of line", desc: "D"},
		item{title: "Delete to end of word", desc: "de"},
		item{title: "Delete to beginning of line", desc: "d0"},
		item{title: "Delete to beginning of word", desc: "db"},
		item{title: "Change in word", desc: "ciw"},
		item{title: "Change in quotes", desc: "ci\""},
		item{title: "Change to end of line", desc: "C"},
		item{title: "Change line", desc: "cc"},
		item{title: "Indent line", desc: ">>/<<"},
		item{title: "Undo", desc: "u"},
		item{title: "Redo", desc: "Ctrl+r"},
		item{title: "Copy/yank", desc: "y"},
		item{title: "Copy/yank line", desc: "yy"},
		item{title: "Paste", desc: "p"},
		item{title: "Paste before cursor", desc: "P"},
		item{title: "Delete", desc: "d"},
		item{title: "Delete line", desc: "dd"},
		item{title: "Delete word", desc: "dw"},
		item{title: "Delete to end of line", desc: "D"},
		item{title: "Delete to end of word", desc: "de"},
		item{title: "Delete to beginning of line", desc: "d0"},
		item{title: "Delete to beginning of word", desc: "db"},
		item{title: "Visual line mode", desc: "V"},
		item{title: "Flip cursor in visual mode", desc: "o"},
		item{title: "Select similar words and replace", desc: "*"},
		item{title: "Search for word", desc: "/word + n/N"},
		item{title: "Search for word backwards", desc: "?word + n/N"},
		item{title: "Search for next word", desc: "n"},
		item{title: "Search for previous word", desc: "N"},
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Vim Motions for Babies"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
