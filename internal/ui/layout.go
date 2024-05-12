package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/umtdemr/spor-istanbul-cli/internal/service"
	"github.com/umtdemr/spor-istanbul-cli/internal/session"
	"log"
)

type screen int

const (
	authScreen screen = iota
	subscriptionScreen
	sessionScreen
	alarmScreen
)

type screenDoneMsg struct{}

func screenDone() tea.Msg {
	return screenDoneMsg{}
}

type model struct {
	currentScreen      screen
	authModel          AuthModel
	subscriptionModel  SubscriptionModel
	sessionScreenModel SessionModel
	alarmScreenModel   AlarmModel
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
		switch m.currentScreen {
		case authScreen:
			m.currentScreen = subscriptionScreen
			return m, m.subscriptionModel.InitSubscriptions()
		case subscriptionScreen:
			m.currentScreen = sessionScreen
			m.sessionScreenModel.selectedSubscriptionId = m.
				subscriptionModel.
				subscriptions[m.subscriptionModel.selectedSubscription].
				PostRequestId
			return m, m.sessionScreenModel.InitSessions()
		case sessionScreen:
			m.currentScreen = alarmScreen

			current := 0

			for _, collection := range m.sessionScreenModel.collections {
				for _, singleSession := range collection.Sessions {
					if current == m.sessionScreenModel.selectedSession {
						m.alarmScreenModel.selectedSession = &session.SelectedSession{
							Day:  collection.Day,
							Date: collection.Date,
							Time: singleSession.Time,
							Id:   singleSession.Id,
						}
						break
					}
					current++
				}
			}

			return m, nil
		}
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
	case sessionScreen:
		newModel, cmd := m.sessionScreenModel.Update(msg)
		m.sessionScreenModel = newModel.(SessionModel)
		return m, cmd
	case alarmScreen:
		newModel, cmd := m.alarmScreenModel.Update(msg)
		m.alarmScreenModel = newModel.(AlarmModel)
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
	case sessionScreen:
		return m.sessionScreenModel.View()
	case alarmScreen:
		return m.alarmScreenModel.View()
	}
	return ""
}

func StartApp() {
	api := service.NewService()
	p := tea.NewProgram(model{
		currentScreen:      authScreen,
		authModel:          initialAuthModel(api),
		subscriptionModel:  initialSubscriptionModel(api),
		sessionScreenModel: initialSessionModel(api),
		alarmScreenModel:   initialAlarmModel(api),
	})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
