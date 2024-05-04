package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/umtdemr/spor-istanbul-cli/internal/session"
	"strings"
)

func GenerateSubscriptionListScreen() string {
	var selectedSubscription int = 0

	doc := strings.Builder{}

	sub1 := session.Subscription{
		Name:      "HAMZA YERLİKAYA SPOR KOMPLEKSİ FİTNESS TÜM GÜN (07:00 22:00) 18 KONTÖR",
		Date:      "15.04.2024 - 14.06.2024",
		Remaining: "9",
	}

	rows := [][]string{
		{sub1.Name, sub1.Date, sub1.Remaining},
		{sub1.Name, sub1.Date, sub1.Remaining},
	}

	subscriptionTable := table.New().Border(lipgloss.RoundedBorder()).
		Headers("Place", "Subscription Date", "Remaining").
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			var style lipgloss.Style

			if row == 0 {
				return lipgloss.NewStyle().AlignHorizontal(lipgloss.Center)
			}
			style = style.AlignHorizontal(lipgloss.Center).Padding(1)

			if row-1 == selectedSubscription {
				style = style.Background(lipgloss.Color("#7239EA")).Foreground(lipgloss.Color("#fff"))
			}

			return style
		})
	doc.WriteString(subscriptionTable.Render())
	return doc.String()
}
