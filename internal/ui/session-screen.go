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
	loading                bool
	err                    error
}

func initialSessionModel(api *service.Service) SessionModel {
	return SessionModel{
		api:     api,
		loading: true,
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
		return m, nil
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

	sessionRenderer := lipgloss.NewStyle().PaddingTop(2).PaddingBottom(2).Width(20).Border(lipgloss.RoundedBorder(), true).MarginRight(1)
	sessionTitle := lipgloss.NewStyle().AlignHorizontal(lipgloss.Center).Width(20)
	panel := lipgloss.NewStyle().AlignHorizontal(lipgloss.Center).Width(20)
	var renderedSessionColumns []string
	for _, sessionList := range collections {
		renderedPanelStr := lipgloss.JoinVertical(
			lipgloss.Top,
			panel.Render(sessionList.Day),
			panel.Render("03.05.2024"),
		)
		renderedSessionRows := []string{renderedPanelStr}
		for _, singleSession := range sessionList.Sessions {
			sessionTitle.MarginTop(0)

			if singleSession.Applicable {
				sessionRenderer.BorderForeground(lipgloss.Color("#00ff00"))
			} else {
				sessionRenderer.BorderForeground(lipgloss.Color("#ff0000"))
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
