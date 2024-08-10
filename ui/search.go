package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/slh335/mc-manager/services"
)

type searchModel struct {
	textinput textinput.Model
	err       error
	mods      []services.ModrinthMod
	loader    string
	mcVersion string
	selected  map[int]bool
	cursor    int
	width     int
	height    int
}

func InitialSearchModel(query, loader, mcVersion string) searchModel {
	ti := textinput.New()
	ti.SetValue(query)
	ti.Placeholder = "Search mods..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 40

	return searchModel{
		textinput: ti,
		err:       nil,
		mods:      []services.ModrinthMod{},
		loader:    loader,
		mcVersion: mcVersion,
		selected:  map[int]bool{},
		cursor:    0,
	}
}

func (m searchModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m searchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.cursor == 0 {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "esc":
				return m, tea.Quit
			case "enter":
				if strings.TrimSpace(m.textinput.Value()) != "" {
					mods, err := services.SearchModrinthMods(m.textinput.Value(), m.loader, m.mcVersion)
					if err != nil {
						m.err = err
					} else {
						m.mods = mods[:min(len(mods), m.height/2)]
						m.cursor = 1
					}
					return m, nil
				}
			}
		case tea.WindowSizeMsg:
			m.height = msg.Height
			m.width = msg.Width
		case error:
			m.err = msg
			return m, nil
		}
		m.textinput, cmd = m.textinput.Update(msg)
		return m, cmd
	} else {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "esc":
				m.cursor = 0
			case "up", "k":
				if m.cursor > 1 {
					m.cursor--
				}
			case "down", "j":
				if len(m.mods) > m.cursor {
					m.cursor++
				}
			case "enter", " ":
				m.selected[m.cursor-1] = !m.selected[m.cursor-1]
			case "i":
				selectedMods := []services.ModrinthMod{}
				for i, selected := range m.selected {
					if selected {
						selectedMods = append(selectedMods, m.mods[i])
					}
				}
				return m, tea.Quit
			case "s":
				m.cursor = 0
			}
		}

	}
	return m, nil
}

func (m searchModel) View() string {
	var mainStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		PaddingLeft(1).
		PaddingRight(2)

	var boldFont = lipgloss.NewStyle().
		Bold(true)

	var grayFont = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#777777"))

	var greenFont = lipgloss.NewStyle().
		Foreground(lipgloss.Color("2"))

	if m.err != nil {
		return m.err.Error()
	}

	s := "Search: "
	if m.cursor == 0 {
		s += m.textinput.View()
	} else {
		s += "  " + m.textinput.Value() + grayFont.Render(" (Press 's' to search)")
	}
	s += "\n\n"

	for i, mod := range m.mods {
		var entry string
		if m.cursor == i+1 {
			entry += ">"
		} else {
			entry += " "
		}
		if m.selected[i] {
			entry += "* "
		} else {
			entry += "  "
		}
		entry += mod.Title
		if m.selected[i] {
			entry = boldFont.Render(entry)
		}
		if m.cursor == i+1 {
			entry = greenFont.Render(entry)
		}
		s += entry + "\n"
	}
	return mainStyle.Render(s)
}
