package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/freecode/freecode/internal/args"
	"github.com/freecode/freecode/internal/style"
)

func Run(a args.Args) error {
	if a.Renderer != "" {
		style.SetDefault(style.Parse(a.Renderer))
	}

	model := NewModel(a)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to start TUI: %w", err)
	}
	return nil
}

func RunHeadless(a args.Args) error {
	fmt.Println("Headless mode is not yet implemented")
	return nil
}