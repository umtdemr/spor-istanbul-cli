package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"time"
)

type model struct {
	view            string
	choices         []time.Time
	sessionChoices  []string
	cursor          int
	selectedDate    time.Time
	selectedSession string
}

func initialModel() model {
	return model{
		// Our to-do list is a grocery list
		choices:        []time.Time{time.Now(), time.Now().Add(time.Hour * 24), time.Now().Add(time.Hour * 48)},
		sessionChoices: []string{"8:30 - 14.30", "17.00 - 22.00"},
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.view == "sessions" {
		return m.UpdateSessions(msg)
	}
	return m.UpdateDates(msg)
}

func (m model) View() string {
	if m.view == "sessions" {
		return m.SessionSelectView()
	} else if m.view == "wait" {
		return m.WaitView()
	}
	return m.DateSelectView()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
