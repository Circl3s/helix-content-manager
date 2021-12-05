package main

import (
	"fmt"
	"os"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/key"
)

type byName []Entry

func (s byName) Len() int {
	return len(s)
}

func (s byName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byName) Less(i, j int) bool {
	return (s[i].Title[0] < s[j].Title[0])
}

var fg_primary 		= lipgloss.Color("#db2b39")
var fg_secondary 	= lipgloss.Color("#777777") // #2e324e
var bg_primary 		= lipgloss.Color("#0f111a")

var container_style = lipgloss.NewStyle().
	Margin(1).
	Padding(1).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(fg_primary)


var entry_style = lipgloss.NewStyle().
	Padding(1)

var entry_title_style = lipgloss.NewStyle().
	Bold(true)

var entry_key_style = lipgloss.NewStyle().
	Foreground(fg_secondary)

var entry_selected_title_style = entry_title_style.Copy().
	Foreground(fg_primary)

type Model struct {
	Index		*Index
	Items 		[]Entry
	Selected 	int
	Active		Entry
}

type KeyMap struct {
	Up 		key.Binding
	Down 	key.Binding
	Enter	key.Binding
	Back	key.Binding
	Quit	key.Binding
}

var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("return"),
		key.WithHelp("enter", "confirm"),
	),
	Back: key.NewBinding(
		key.WithKeys("backspace", "escape"),
		key.WithHelp("backspace/esc", "go back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, DefaultKeyMap.Up):
			if (m.Selected != 0) {
				m.Selected -= 1
			}
		case key.Matches(msg, DefaultKeyMap.Down):
			if (m.Selected != len(m.Items) - 1) {
				m.Selected += 1
			}
		}
	case tea.WindowSizeMsg:
		container_style.
			Width(msg.Width - 4).
			Height(msg.Height - 4)
	}

	var cmd tea.Cmd

	return m, cmd
}

func (m Model) View() string {
	var s string
	var title string
	for i, e := range m.Items {
		if i == m.Selected {
			title = entry_selected_title_style.Render(e.Title)
		} else {
			title = entry_title_style.Render(e.Title)
		}
		
		key := entry_key_style.Render(e.Key)

		s += entry_style.Render(lipgloss.JoinVertical(lipgloss.Left, title, key))

		
	}
	return container_style.Render(s)
}

func NewModel(index *Index) Model {
	var items []Entry
	for _, v := range index.Entries {
		items = append(items, v)
	}
	sort.Sort(byName(items))
	return Model{Index: index, Items: items, Selected: 0, Active: Entry{}}
}

func (index *Index) TUI() {
	var entries []Entry
	for _, v := range index.Entries {
		entries = append(entries, v)
	}
	m := NewModel(index)
	
	p := tea.NewProgram(m)
	p.EnterAltScreen()

	err := p.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}