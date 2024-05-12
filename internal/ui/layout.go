package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/umtdemr/spor-istanbul-cli/internal/service"
	"log"
)

type screen int

const (
	authScreen screen = iota
	subscriptionScreen
	sessionScreen
)

type screenDoneMsg struct{}

func screenDone() tea.Msg {
	return screenDoneMsg{}
}

type model struct {
	currentScreen     screen
	authModel         AuthModel
	subscriptionModel SubscriptionModel
}

func (m model) Init() tea.Cmd {
	// Initialize the first screen
	return m.authModel.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case screenDoneMsg:
		m.currentScreen = subscriptionScreen
		return m, m.subscriptionModel.InitSubscriptions()
	}
	switch m.currentScreen {
	case authScreen:
		newModel, cmd := m.authModel.Update(msg)
		m.authModel = newModel.(AuthModel)
		return m, cmd
	case subscriptionScreen:
		newModel, cmd := m.subscriptionModel.Update(msg)
		m.subscriptionModel = newModel.(SubscriptionModel)
		return m, cmd

	}

	return m, nil
}

func (m model) View() string {
	switch m.currentScreen {
	case authScreen:
		return m.authModel.View()
	case subscriptionScreen:
		return m.subscriptionModel.View()
	}
	return ""
}

func StartApp() {
	api := service.NewService()
	p := tea.NewProgram(model{
		currentScreen:     authScreen,
		authModel:         initialAuthModel(api),
		subscriptionModel: initialSubscriptionModel(api),
	})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
