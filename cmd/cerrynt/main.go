package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/mako/cerrynt-cli/internal/app"
)

func main() {
	p := tea.NewProgram(
		app.New(),
		tea.WithAltScreen(), // Use the terminal's alternate screen buffer.
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "cerrynt: %v\n", err)
		os.Exit(1)
	}
}
