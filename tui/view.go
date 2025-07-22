package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if m.tab == 0 {
		return m.renderNetworkTab()
	}
	return m.renderProcessTab()
}

func (m model) renderNetworkTab() string {
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true).Align(lipgloss.Center).Width(m.width - 10).MarginBottom(1)
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#A0A0A0")).MarginBottom(1)
	statStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Margin(0, 1, 1, 1)
	graphTitleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Margin(1, 0, 0, 0)
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).MarginTop(1)

	headerTabs := ""
	for i, name := range m.interfaces {
		if i == m.selected {
			headerTabs += lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#7D56F4")).Padding(0, 1).Render(" " + name + " ")
		} else {
			headerTabs += lipgloss.NewStyle().Foreground(lipgloss.Color("#A0A0A0")).Padding(0, 1).Render(" " + name + " ")
		}
	}

	footerTabs := ""
	bottomTabs := []string{"Network", "Process"}
	for i, name := range bottomTabs {
		if i == m.tab {
			footerTabs += lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#7D56F4")).Padding(0, 1).Render(" " + name + " ")
		} else {
			footerTabs += lipgloss.NewStyle().Foreground(lipgloss.Color("#A0A0A0")).Padding(0, 1).Render(" " + name + " ")
		}
	}

	stats := lipgloss.JoinHorizontal(
		lipgloss.Center,
		statStyle.Copy().Foreground(lipgloss.Color("#6ADB78")).Render("⬇ "+formatBitsPerSecond(m.rxBps)),
		statStyle.Copy().Foreground(lipgloss.Color("#DB6A6A")).Render("⬆ "+formatBitsPerSecond(m.txBps)),
	)

	graph := func(data []float64, color string) string {
		if len(data) == 0 {
			return ""
		}

		maxVal := 1.0
		for _, v := range data {
			if v > maxVal {
				maxVal = v
			}
		}

		graphRows := make([]string, graphHeight)
		for i := range graphRows {
			graphRows[i] = strings.Repeat(" ", len(data))
		}

		for i, val := range data {
			height := int((val / maxVal) * float64(graphHeight))
			if height > graphHeight {
				height = graphHeight
			}
			for j := 0; j < height; j++ {
				row := []rune(graphRows[graphHeight-1-j])
				if i < len(row) {
					row[i] = '▀'
				}
				graphRows[graphHeight-1-j] = string(row)
			}
		}

		return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(strings.Join(graphRows, "\n"))
	}

	rxGraph := graphTitleStyle.Render("Download") + "\n" + graph(m.rxHistory, "#6ADB78")
	txGraph := graphTitleStyle.Render("Upload") + "\n" + graph(m.txHistory, "#DB6A6A")

	graphs := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(m.width/3-1).Render(rxGraph),
		lipgloss.NewStyle().Width(m.width/2-1).Render(txGraph),
	)

	help := helpStyle.Render("←/→: switch interface  •  tab/1/2: switch view  •  q/esc: quit")

	content := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("Bitnix Network Monitor"),
		headerTabs,
		infoStyle.Render("Interface: "+m.interfaces[m.selected]),
		stats,
		graphs,
		help,
		footerTabs,
	)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Top,
		lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#444444")).Padding(1, 2).Render(content),
	)
}

func (m model) renderProcessTab() string {
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true).Align(lipgloss.Center).Width(m.width - 10).MarginBottom(1)
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).MarginTop(1)

	footerTabs := ""
	bottomTabs := []string{"Network", "Process"}
	for i, name := range bottomTabs {
		if i == m.tab {
			footerTabs += lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#7D56F4")).Padding(0, 1).Render(" " + name + " ")
		} else {
			footerTabs += lipgloss.NewStyle().Foreground(lipgloss.Color("#A0A0A0")).Padding(0, 1).Render(" " + name + " ")
		}
	}

	help := helpStyle.Render("r: sort by rx  •  t: sort by tx  •  b: sort by total")
	content := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("Bitnix Process Monitor"),
		lipgloss.NewStyle().MarginBottom(1).Render(m.table.View()),
		help,
		footerTabs,
	)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Top,
		lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#444444")).Padding(1, 2).Render(content),
	)
}