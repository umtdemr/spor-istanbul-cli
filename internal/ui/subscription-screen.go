package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/umtdemr/spor-istanbul-cli/internal/service"
	"github.com/umtdemr/spor-istanbul-cli/internal/session"
	"strings"
)

type getSubscriptionsMsg struct{}

func getSubscriptions() tea.Msg {
	return getSubscriptionsMsg{}
}

type SubscriptionModel struct {
	api                  *service.Service
	subscriptions        []*session.Subscription
	selectedSubscription int
	loading              bool
	err                  error
}

func initialSubscriptionModel(api *service.Service) SubscriptionModel {
	return SubscriptionModel{
		api:                  api,
		selectedSubscription: -1,
		loading:              true,
	}
}

func (m SubscriptionModel) callSubscriptionsApiCmd() tea.Cmd {
	return func() tea.Msg {
		subscriptions := m.api.GetSubscriptions()
		return subscriptions
	}
}

func (m SubscriptionModel) Init() tea.Cmd {
	return nil
}

func (m SubscriptionModel) InitSubscriptions() tea.Cmd {
	return m.callSubscriptionsApiCmd()
}

func (m SubscriptionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case []*session.Subscription:
		m.subscriptions = msg
		m.loading = false
		m.selectedSubscription = 0
		return m, nil
	case error:
		m.err = msg
		m.loading = false
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.selectedSubscription >= 1 {
				m.selectedSubscription -= 1
			} else {
				m.selectedSubscription = len(m.subscriptions) - 1
			}
			return m, nil
		case tea.KeyDown:
			if m.selectedSubscription < len(m.subscriptions)-1 {
				m.selectedSubscription += 1
			} else {
				m.selectedSubscription = 0
			}
			return m, nil
		case tea.KeyEnter:
			return m, screenDone
		}

	}
	return m, nil
}

func (m SubscriptionModel) View() string {
	if m.loading {
		return "loading...."
	}
	return m.GenerateSubscriptionListScreen(m.subscriptions)
}

func (m SubscriptionModel) GenerateSubscriptionListScreen(subscriptions []*session.Subscription) string {

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

			if row-1 == m.selectedSubscription {
				style = style.Background(lipgloss.Color("#7239EA")).Foreground(lipgloss.Color("#fff"))
			}

			return style
		})
	doc.WriteString(subscriptionTable.Render())
	return doc.String()
}
