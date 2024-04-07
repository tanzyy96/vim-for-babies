package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tanzyy96/vim-for-babies/db"
)

var (
	docStyle      = lipgloss.NewStyle().Margin(1, 2)
	addTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type viewMode string

const (
	viewModeList   viewMode = "list"
	viewModeAdd    viewMode = "add"
	viewModeDelete viewMode = "delete"
)

// Core model for storing states in BubbleTea
type model struct {
	db       db.DB
	list     list.Model
	viewMode viewMode

	inputs     []textinput.Model
	focusIndex int

	addMotionHelp       help.Model
	addMotionHelpKeymap additionalKeymap

	deleteSelectedIndex int
	deleteInput         textinput.Model
	// deleteHelp  help.Model
	// deleteHelpKeymap additionalKeymap
}

// Interface for new keymaps
type additionalKeymap interface {
	Bindings() []key.Binding
}

// Get new textInputModels for AddMotion. Use this function to reset the inputs
func textInputModel() []textinput.Model {
	inputs := make([]textinput.Model, 2)

	var t textinput.Model
	for i := range inputs {
		t = textinput.New()
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Description of Motion"
			t.Focus()
		case 1:
			t.Placeholder = "Motion Keys"
			t.CharLimit = 64
		}
		inputs[i] = t
	}

	return inputs
}

func newDeleteTextInput() textinput.Model {
	t := textinput.New()
	t.Placeholder = "y/n"
	t.CharLimit = 5
	return t
}

func initialModel(db db.DB, items []list.Item) model {
	listKeys := newListKeymap()
	addMotionKeys := newAddMotionKeymap()

	m := model{
		db:                  db,
		inputs:              textInputModel(),
		list:                list.New(items, list.NewDefaultDelegate(), 0, 0),
		addMotionHelp:       help.New(),
		addMotionHelpKeymap: addMotionKeys,
		deleteInput:         newDeleteTextInput(),
	}

	m.list.Title = "Vim Motions for Babies"
	m.list.AdditionalShortHelpKeys = listKeys.Bindings
	m.list.AdditionalFullHelpKeys = listKeys.Bindings

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// Main update function which handles inputs and key presses
// Calls other update functions based on the viewMode
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.viewMode {
	case viewModeAdd:
		return m.updateAddMotion(msg)
	case viewModeDelete:
		return m.deleteUpdate(msg)
	default:
		return m.updateList(msg)
	}
}

func (m *model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.IsFiltered() {
			if msg.String() == "esc" {
				return m, tea.Quit
			}
		}

		if msg.String() == "n" {
			for _, input := range m.inputs {
				input.SetValue("")
			}
			m.viewMode = viewModeAdd
			return m, nil
		}
		if msg.String() == "d" {
			m.viewMode = viewModeDelete
			m.deleteInput.Focus()
			m.deleteInput.SetValue("")
			m.deleteSelectedIndex = m.list.Index()
			return m, nil
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *model) updateAddMotion(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.viewMode = viewModeList
			return m, nil

		case "enter":
			// Save the motion
			desc := m.inputs[0].Value()
			keys := m.inputs[1].Value()

			if desc == "" || keys == "" {
				m.viewMode = viewModeList
				return m, nil
			}

			newCommand := db.Command{Title: desc, Description: keys}
			if err := m.db.Add(newCommand); err != nil {
				fmt.Println("Error adding command:", err)
				// temporarily force quit the app
				return m, tea.Quit
			}

			m.list.InsertItem(-1, item{title: desc, desc: keys})
			m.viewMode = viewModeList

			return m, nil

		case "up", "down":
			if msg.String() == "up" {
				if m.focusIndex > 0 {
					m.focusIndex--
				}
			}
			if msg.String() == "down" {
				if m.focusIndex < len(m.inputs)-1 {
					m.focusIndex++
				}
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m model) deleteUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.viewMode = viewModeList
			return m, nil
		case "enter":
			validConfirms := []string{"y", "yes"}
			if m.deleteInput.Value() == "" {
				m.viewMode = viewModeList
				return m, nil
			} else if !slices.Contains(validConfirms, strings.ToLower(m.deleteInput.Value())) {
				m.deleteInput.SetValue("")
				m.viewMode = viewModeList
				return m, nil
			}

			selected := m.list.SelectedItem().(item)
			selectedIndex := m.list.Index()
			if err := m.db.Delete(selected.title); err != nil {
				fmt.Println("Error deleting command:", err)
				// temporarily force quit the app
				return m, tea.Quit
			}
			m.list.RemoveItem(selectedIndex)
			m.viewMode = viewModeList
			return m, nil
		}
	}

	deleteModel, cmd := m.deleteInput.Update(msg)
	m.deleteInput = deleteModel

	return m, cmd
}

func (m model) View() string {
	switch m.viewMode {
	case viewModeAdd:
		return m.addMotionView()
	case viewModeDelete:
		return m.deleteView()
	default:
		return m.listView()
	}
}

func (m model) addMotionView() string {
	var b strings.Builder
	b.WriteString(lipgloss.NewStyle().Render("Add a new motion\n"))

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteString("\n")
		}
	}

	helpView := m.addMotionHelp.View(m.addMotionHelpKeymap.(help.KeyMap))
	b.WriteString("\n\n\n" + helpView)

	return b.String()
}

func (m model) listView() string {
	return docStyle.Render(m.list.View())
}

func (m model) deleteView() string {
	var b strings.Builder
	b.WriteString(lipgloss.NewStyle().Render("Are you sure you want to delete this motion? (y/n)\n\n"))
	// TODO: add some styling for the selected item title
	b.WriteString(m.list.SelectedItem().(item).title + "\n")

	b.WriteString(m.deleteInput.View() + "\n")

	helpView := m.addMotionHelp.View(m.addMotionHelpKeymap.(help.KeyMap))
	b.WriteString("\n\n\n" + helpView)

	return b.String()
}

// type errMsg struct {
// 	err error
// }

func main() {
	commandDb := db.New()
	items := []list.Item{}
	commands, err := commandDb.GetAll()
	if err != nil {
		fmt.Println("Error getting commands:", err)
		os.Exit(1)
	}
	for _, command := range commands {
		items = append(items, item{title: command.Title, desc: command.Description})
	}

	p := tea.NewProgram(initialModel(commandDb, items), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
