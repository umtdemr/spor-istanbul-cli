package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

func LoginScreen() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type model struct {
	textInput textinput.Model
	username  string
	password  string
	mode      string
	err       error
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "username"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		err:       nil,
		mode:      "username",
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.mode == "username" {
				m.username = m.textInput.Value()
				m.mode = "password"
				m.textInput.Placeholder = "password"
				m.textInput.EchoMode = textinput.EchoPassword
				m.textInput.SetValue("")
				return m, tea.ClearScrollArea
			}
			m.password = m.textInput.Value()
			return m, tea.Quit

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	title := "username"
	if m.mode != "username" {
		title = "password"
	}
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		title,
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}
