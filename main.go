package main

import (
	"fmt"
	"os"

	"bitnix/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	fmt.Println("Bitnix — Real-time Network Monitor")

	p := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}
