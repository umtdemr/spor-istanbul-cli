package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/umtdemr/spor-istanbul-cli/internal/service"
	"github.com/umtdemr/spor-istanbul-cli/internal/session"
	"golang.org/x/term"
	"os"
	"strings"
)

const (
	width = 96

	columnWidth = 30
)

var (
	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)

type SessionModel struct {
	api                    *service.Service
	collections            []*session.Collection
	selectedSubscriptionId string
	selectedSession        int
	totalSessionLength     int
	loading                bool
	err                    error
}

func initialSessionModel(api *service.Service) SessionModel {
	return SessionModel{
		api:             api,
		loading:         true,
		selectedSession: -1,
	}
}

func (m SessionModel) callSessionsApiCmd() tea.Cmd {
	return func() tea.Msg {
		collections := m.api.GetSessions(m.selectedSubscriptionId)
		return collections
	}
}

func (m SessionModel) Init() tea.Cmd {
	return nil
}

func (m SessionModel) InitSessions() tea.Cmd {
	return m.callSessionsApiCmd()
}

func (m SessionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case []*session.Collection:
		m.collections = msg
		m.loading = false
		m.selectedSession = 0

		totalLength := 0

		for _, collection := range msg {
			totalLength += len(collection.Sessions)
		}

		m.totalSessionLength = totalLength

		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.selectedSession >= 1 {
				m.selectedSession -= 1
			} else {
				m.selectedSession = m.totalSessionLength - 1
			}
			return m, nil
		case tea.KeyDown:
			if m.selectedSession < m.totalSessionLength-1 {
				m.selectedSession += 1
			} else {
				m.selectedSession = 0
			}
			return m, nil
		}
	}
	return m, nil
}

func (m SessionModel) View() string {
	if m.loading {
		return "loading"
	}
	return m.GenerateSessionScreen(m.collections)
}

func (m SessionModel) GenerateSessionScreen(collections []*session.Collection) string {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	doc := strings.Builder{}

	{
		title := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("Please select a session")

		dialog := lipgloss.Place(width, 3,
			lipgloss.Center, lipgloss.Center,
			dialogBoxStyle.Render(title),
			lipgloss.WithWhitespaceChars("-"),
			lipgloss.WithWhitespaceForeground(subtle),
		)

		doc.WriteString(dialog + "\n\n")
	}

	var renderedSessionColumns []string
	currentSession := 0
	for _, sessionList := range collections {
		panel := lipgloss.NewStyle().AlignHorizontal(lipgloss.Center).Width(20)
		renderedPanelStr := lipgloss.JoinVertical(
			lipgloss.Top,
			panel.Render(sessionList.Day),
			panel.Render("03.05.2024"),
		)
		renderedSessionRows := []string{renderedPanelStr}
		for _, singleSession := range sessionList.Sessions {
			sessionTitle := lipgloss.NewStyle().AlignHorizontal(lipgloss.Center).Width(20)
			sessionRenderer := lipgloss.NewStyle().PaddingTop(2).PaddingBottom(2).Width(20).Border(lipgloss.RoundedBorder(), true).MarginRight(1)
			if singleSession.Applicable {
				sessionRenderer.BorderForeground(lipgloss.Color("#00ff00"))
			} else {
				sessionRenderer.BorderForeground(lipgloss.Color("#ff0000"))
			}

			if currentSession == m.selectedSession {
				sessionRenderer.Background(lipgloss.Color("#7239EA"))
				sessionRenderer.Foreground(lipgloss.Color("#FFF"))
			}

			applicableText := "Yer Var"
			if !singleSession.Applicable {
				applicableText = "Dolu"
			}

			details := lipgloss.JoinVertical(
				lipgloss.Top,
				sessionTitle.Render(
					fmt.Sprintf("%s / %s", singleSession.Available, singleSession.Limit),
				),
				sessionTitle.MarginTop(1).Render(singleSession.Time),
				sessionTitle.MarginTop(1).Render(applicableText),
			)
			renderedSessionRows = append(renderedSessionRows, sessionRenderer.Render(details))
			currentSession++
		}

		renderedSessionColumns = append(
			renderedSessionColumns,
			lipgloss.JoinVertical(
				lipgloss.Left,
				renderedSessionRows...,
			),
		)
	}

	allRendered := lipgloss.JoinHorizontal(lipgloss.Left, renderedSessionColumns...)
	doc.WriteString(allRendered)

	if physicalWidth > 0 {
		docStyle = docStyle.MaxWidth(physicalWidth)
	}

	return docStyle.Render(doc.String())
}
