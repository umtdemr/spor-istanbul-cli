package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/umtdemr/spor-istanbul-cli/internal/service"
	"github.com/umtdemr/spor-istanbul-cli/internal/session"
)

type AlarmModel struct {
	api             *service.Service
	selectedSession *session.SelectedSession
	err             error
}

func initialAlarmModel(api *service.Service) AlarmModel {
	return AlarmModel{
		api: api,
	}
}

func (m AlarmModel) Init() tea.Cmd {
	return nil
}

func (m AlarmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m AlarmModel) View() string {
	return fmt.Sprintf("%s %s - %s\n", m.selectedSession.Day, m.selectedSession.Date, m.selectedSession.Time)
}
