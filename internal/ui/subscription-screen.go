package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/umtdemr/spor-istanbul-cli/internal/session"
	"strings"
)

func GenerateSubscriptionListScreen(subscriptions []*session.Subscription) string {
	var selectedSubscription int = 0

	doc := strings.Builder{}

	rows := make([][]string, len(subscriptions))

	for i, subscription := range subscriptions {
		thisRow := []string{subscription.Name, subscription.Date, subscription.Remaining}
		rows[i] = thisRow
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
