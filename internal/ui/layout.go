package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/umtdemr/spor-istanbul-cli/internal/service"
	"github.com/umtdemr/spor-istanbul-cli/internal/session"
	"golang.org/x/term"
	"log"
	"os"
	"strings"
)

type screen int

const (
	authScreen screen = iota
	subscriptionScreen
	sessionScreen
	alarmScreen
)

var terminalWidth, terminalHeight, _ = term.GetSize(int(os.Stdout.Fd()))

var (
	ContainerStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Border(lipgloss.NormalBorder())
	titleStyle     = lipgloss.NewStyle().Width(terminalWidth - 15).Align(lipgloss.Center)
	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)
)

// screenDoneMsg is a type that tells Update method to change the screen
type screenDoneMsg struct{}

// screenDone returns screenDoneMsg
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
	case screenDoneMsg: // cmd to change the page
		switch m.currentScreen {
		case authScreen:
			m.currentScreen = subscriptionScreen
			return m, m.subscriptionModel.InitSubscriptions()
		case subscriptionScreen:
			m.currentScreen = sessionScreen
			// set selected subscription for session model
			selectedSubscriptionId := m.
				subscriptionModel.
				subscriptions[m.subscriptionModel.selectedSubscription].
				PostRequestId

			m.sessionScreenModel.selectedSubscriptionId = selectedSubscriptionId
			m.alarmScreenModel.selectedSubscriptionId = selectedSubscriptionId
			return m, m.sessionScreenModel.InitSessions()
		case sessionScreen:
			m.currentScreen = alarmScreen

			current := 0

			// set selected session for alarm model
			found := false
			for _, collection := range m.sessionScreenModel.collections {
				if found {
					break
				}
				for _, singleSession := range collection.Sessions {
					if current == m.sessionScreenModel.selectedSession {
						m.alarmScreenModel.selectedSession = &session.SelectedSession{
							Day:  collection.Day,
							Date: collection.Date,
							Time: singleSession.Time,
							Id:   singleSession.Id,
						}
						found = true
						break
					}
					current++
				}
			}

			return m, m.alarmScreenModel.InitAlarm()
		}
	}

	// otherwise, handle the cmd with the current screen's models
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

// View is the main layout for the CLI
func (m model) View() string {
	// get the current screen view
	view := ""

	switch m.currentScreen {
	case authScreen:
		view = m.authModel.View()
	case subscriptionScreen:
		view = m.subscriptionModel.View()
	case sessionScreen:
		view = m.sessionScreenModel.View()
	case alarmScreen:
		view = m.alarmScreenModel.View()
	}

	doc := strings.Builder{}

	// render the screen in the container
	mainView := ContainerStyle.Width(terminalWidth - 2).Height(terminalHeight - 4).MaxHeight(terminalHeight).Render(view)

	// add footer
	footer := lipgloss.NewStyle().
		Width(terminalWidth - 2).
		Align(lipgloss.Center)

	doc.WriteString(mainView)
	doc.WriteString("\n\n")
	doc.WriteString(footer.Render("↑/↓ to select"))
	return doc.String()
}

// StartApp starts the CLI application
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
