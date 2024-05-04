package ui

import (
	"github.com/charmbracelet/lipgloss"
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

func GenerateSessionScreen() string {
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

	sessionsFirst := []session.Session{
		{
			Day:       "Saturday",
			Date:      "03.05.2024",
			Available: 0,
			Limit:     300,
			Time:      "08:00 - 16:30",
		},
		{
			Day:       "Saturday",
			Date:      "03.05.2024",
			Available: 1,
			Limit:     240,
			Time:      "17:00 - 21:00",
		},
	}

	sessionsSecond := []session.Session{
		{
			Day:       "Saturday",
			Date:      "03.05.2024",
			Available: 1,
			Limit:     300,
			Time:      "08:00 - 16:30",
		},
		{
			Day:       "Saturday",
			Date:      "03.05.2024",
			Available: 0,
			Limit:     240,
			Time:      "17:00 - 21:00",
		},
	}

	sessionsThird := []session.Session{
		{
			Day:       "Saturday",
			Date:      "03.05.2024",
			Available: 1,
			Limit:     300,
			Time:      "08:00 - 16:30",
		},
		{
			Day:       "Saturday",
			Date:      "03.05.2024",
			Available: 1,
			Limit:     240,
			Time:      "17:00 - 21:00",
		},
	}
	sessionsLast := []session.Session{
		{
			Day:       "Saturday",
			Date:      "03.05.2024",
			Available: 0,
			Limit:     300,
			Time:      "08:00 - 16:30",
		},
		{
			Day:       "Saturday",
			Date:      "03.05.2024",
			Available: 0,
			Limit:     240,
			Time:      "17:00 - 21:00",
		},
	}

	allSessions := [][]session.Session{
		sessionsFirst,
		sessionsSecond,
		sessionsThird,
		sessionsLast,
	}
	sessionRenderer := lipgloss.NewStyle().Height(10).Width(20).Border(lipgloss.RoundedBorder(), true).MarginRight(1)
	sessionTitle := lipgloss.NewStyle().AlignHorizontal(lipgloss.Center).Width(20)
	var renderedSessionColumns []string
	for _, sessionList := range allSessions {
		var renderedSessionRows []string
		for _, singleSession := range sessionList {
			sessionTitle.MarginTop(0)

			if singleSession.Available > 0 {
				sessionRenderer.BorderForeground(lipgloss.Color("#00ff00"))
			} else {
				sessionRenderer.BorderForeground(lipgloss.Color("#ff0000"))
			}

			details := lipgloss.JoinVertical(
				lipgloss.Top,
				sessionTitle.Render(singleSession.Day),
				sessionTitle.Render("03.05.2024"),
				sessionTitle.MarginTop(1).Render("239 / 240"),
				sessionTitle.Render("Year Var"),
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
