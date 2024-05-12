package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/umtdemr/spor-istanbul-cli/internal/service"
	"github.com/umtdemr/spor-istanbul-cli/internal/session"
	"time"
)

type AlarmModel struct {
	api                    *service.Service
	selectedSession        *session.SelectedSession
	selectedSubscriptionId string
	checkCount             int
	sub                    chan bool
	err                    error
}

type responseMsg bool

func initialAlarmModel(api *service.Service) AlarmModel {
	return AlarmModel{
		api: api,
		sub: make(chan bool),
	}
}

func (m AlarmModel) Init() tea.Cmd {
	return nil
}

func (m AlarmModel) listenForActivity() tea.Cmd {
	return func() tea.Msg {
		for {
			val := m.api.CheckSessionApplicable(m.selectedSubscriptionId, m.selectedSession.Id)
			if val {
				m.sub <- val
				return nil
			}
			m.sub <- val
			time.Sleep(5 * time.Second)
		}
	}
}

func (m AlarmModel) waitForActivity() tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-m.sub)
	}
}

func (m AlarmModel) alarmCmd() tea.Cmd {
	return tea.Batch(
		m.listenForActivity(),
		m.waitForActivity(),
	)
}

func (m AlarmModel) InitAlarm() tea.Cmd {
	return m.alarmCmd()
}

func (m AlarmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case responseMsg:
		m.checkCount++
		if msg {
			close(m.sub)
			return m, nil
		}
		return m, m.waitForActivity()
	}
	return m, nil
}

func (m AlarmModel) View() string {
	return fmt.Sprintf("%s %s - %s - %v \n", m.selectedSession.Day, m.selectedSession.Date, m.selectedSession.Time, m.checkCount)
}
