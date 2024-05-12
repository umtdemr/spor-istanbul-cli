package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/umtdemr/spor-istanbul-cli/internal/service"
)

type (
	errMsg error
)

type AuthModel struct {
	api       *service.Service
	textInput textinput.Model
	username  string
	password  string
	mode      string
	loggedErr string
	loading   bool
	err       error
}

func initialAuthModel(api *service.Service) AuthModel {
	ti := textinput.New()
	ti.Placeholder = "username"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return AuthModel{
		api:       api,
		textInput: ti,
		err:       nil,
		mode:      "username",
	}
}

func (m AuthModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m AuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

			isLoggedIn := m.api.Login(m.username, m.password)

			if !isLoggedIn {
				m.textInput.SetValue("")
				m.loggedErr = "try again"
				m.mode = "username"
				m.textInput.EchoMode = textinput.EchoNormal
				m.textInput.Placeholder = "username"
				return m, nil
			}

			return m, screenDone
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m AuthModel) View() string {
	title := "username"
	if m.mode != "username" {
		title = "password"
	}

	var screenTitle string

	if m.loggedErr != "" {
		screenTitle = fmt.Sprintf(" - %s", m.loggedErr)
	}

	title += screenTitle

	return fmt.Sprintf(
		"%s\n\n%s",
		title,
		m.textInput.View(),
	) + "\n"
}
