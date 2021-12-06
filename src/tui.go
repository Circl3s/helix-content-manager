package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

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
	BorderStyle(lipgloss.RoundedBorder())


var entry_style = lipgloss.NewStyle().
	Padding(1)

var entry_title_style = lipgloss.NewStyle().
	Bold(true)

var entry_key_style = lipgloss.NewStyle().
	Foreground(fg_secondary)

var entry_selected_title_style = entry_title_style.Copy().
	Foreground(fg_primary)

var active_border = lipgloss.NewStyle().
	BorderForeground(fg_primary)

var inactive_border = lipgloss.NewStyle().
	BorderForeground(fg_secondary)

type Model struct {
	Index		*Index
	Items 		[]Entry
	Selected 	int
	Active		Entry
	Editing		bool
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
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	),
	Back: key.NewBinding(
		key.WithKeys("backspace", "esc"),
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
		if (!m.Editing) {
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
			case key.Matches(msg, DefaultKeyMap.Enter):
				if (!m.Editing) {
					m.Editing = true
				}
			}
		} else {
			switch {
			case key.Matches(msg, DefaultKeyMap.Back):
				m.Editing = false
			}
		}
		
	case tea.WindowSizeMsg:
		container_style.
			Width((msg.Width / 2) - 4).
			Height(msg.Height - 4)
	}

	var cmd tea.Cmd

	return m, cmd
}

func (m Model) View() string {
	var s string
	var r string
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

	selected_entry := m.Index.Entries[m.Items[m.Selected].Key]
	r = selected_entry.Title + "\n" + selected_entry.Key + "\n" + strings.Join(selected_entry.Tags, ", ")
	var l_style, r_style lipgloss.Style
	if (m.Editing) {
		l_style = container_style.Copy().Inherit(inactive_border)
		r_style = container_style.Copy().Inherit(active_border)
	} else {
		l_style = container_style.Copy().Inherit(active_border)
		r_style = container_style.Copy().Inherit(inactive_border)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, l_style.Render(s), r_style.Render(r))
}

func NewModel(index *Index) Model {
	var items []Entry
	for _, v := range index.Entries {
		items = append(items, v)
	}
	sort.Sort(byName(items))
	return Model{Index: index, Items: items, Selected: 0, Active: Entry{}, Editing: false}
}

func (index *Index) TUI() {
	var entries []Entry
	for _, v := range index.Entries {
		entries = append(entries, v)
	}
	m := NewModel(index)
	
	p := tea.NewProgram(m, tea.WithMouseCellMotion())
	p.EnterAltScreen()

	err := p.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}