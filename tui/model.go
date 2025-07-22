package tui

import (
	"fmt"
	"time"

	"bitnix/server"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/table"
)

const (
	historySize = 30
	graphHeight = 8
)


type model struct {
	interfaces    []string
	selected      int
	rxBps         float64
	txBps         float64
	rxHistory     []float64
	txHistory     []float64
	width         int
	height        int
	tab           int // 0 = Network, 1 = Process
	processes     []server.ProcessNetStats
	table         table.Model
	sortBy        SortBy
}

func NewModel() model {
	ifaces := server.ListInterfaces()
	columns := []table.Column{
		{Title: "PID", Width: 6},
		{Title: "NAME", Width: 24},
		{Title: "RX", Width: 12},
		{Title: "TX", Width: 12},
		{Title: "TOTAL", Width: 12},
	}
	tbl := table.New(table.WithColumns(columns), table.WithHeight(15))
	tbl.SetStyles(defaultTableStyle())

	return model{
		interfaces: ifaces,
		selected:   0,
		rxHistory:  make([]float64, 0, historySize),
		txHistory:  make([]float64, 0, historySize),
		tab:        0,
		table:      tbl,
		sortBy:     SortByTotal,
	}
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tick(), tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "left":
			if m.selected > 0 {
				m.selected--
			}
		case "right":
			if m.selected < len(m.interfaces)-1 {
				m.selected++
			}
		case "tab", "1", "2":
			m.tab = (m.tab + 1) % 2
		case "r":
			m.sortBy = SortByRx
		case "t":
			m.sortBy = SortByTx
		case "b":
			m.sortBy = SortByTotal
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetWidth(m.width - 6)

	case tickMsg:
		m.interfaces = server.ListInterfaces()
		if len(m.interfaces) == 0 {
			return m, tick()
		}
		if m.selected >= len(m.interfaces) {
			m.selected = 0
		}

		iface := m.interfaces[m.selected]
		stats := server.GetNetStats(iface)
		m.rxBps = stats.RxBps
		m.txBps = stats.TxBps
		m.rxHistory = append(m.rxHistory, m.rxBps)
		m.txHistory = append(m.txHistory, m.txBps)
		if len(m.rxHistory) > historySize {
			m.rxHistory = m.rxHistory[1:]
		}
		if len(m.txHistory) > historySize {
			m.txHistory = m.txHistory[1:]
		}

		m.processes = GetAllProcessNetStatsSorted(m.sortBy)
		m.table.SetRows(processRows(m.processes))

		return m, tick()
	}

	if m.tab == 1 {
		m.table, cmd = m.table.Update(msg)
	}

	return m, cmd
}

func processRows(procs []server.ProcessNetStats) []table.Row {
	rows := make([]table.Row, 0, len(procs))
	for _, p := range procs {
		total := p.RxBps + p.TxBps
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", p.PID),
			truncate(p.Name, 22),
			formatBitsPerSecond(p.RxBps),
			formatBitsPerSecond(p.TxBps),
			formatBitsPerSecond(total),
		})
	}
	return rows
}

func formatBitsPerSecond(bps float64) string {
	units := []string{"bps", "Kbps", "Mbps", "Gbps"}
	unitIndex := 0
	for bps >= 1000 && unitIndex < len(units)-1 {
		bps /= 1000
		unitIndex++
	}
	return fmt.Sprintf("%.1f %s", bps, units[unitIndex])
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func defaultTableStyle() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Foreground(lipgloss.Color("#ffffff")).
		Bold(true)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#ffffff")).
		Background(lipgloss.Color("#7D56F4"))

	return s
}

