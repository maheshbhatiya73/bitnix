package tui

import (
	"fmt"
	"time"

	"bitnix/server"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	rxMbps float64
	txMbps float64
	quit   bool
}

func NewModel() model {
	return model{}
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			m.quit = true
			return m, tea.Quit
		}
	case tickMsg:
		stats := server.GetNetStats()
		m.rxMbps = stats.RxMbps
		m.txMbps = stats.TxMbps
		return m, tick()
	}
	return m, nil
}

func (m model) View() string {
	if m.quit {
		return "Bitnix closed.\n"
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		Padding(0, 1)

	statStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Bold(true)

	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Margin(1).
		Width(40)

	content := fmt.Sprintf(
		"\n%s\n\n⬇ Download: %s Mbps\n⬆ Upload:   %s Mbps\n\nPress 'q' to quit.",
		titleStyle.Render("Bitnix - Live Network Monitor"),
		statStyle.Render(fmt.Sprintf("%.2f", m.rxMbps)),
		statStyle.Render(fmt.Sprintf("%.2f", m.txMbps)),
	)

	return border.Render(content)
}
