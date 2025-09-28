package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	softPink = lipgloss.Color("#FFB6C1")
	deepPink = lipgloss.Color("#FF1493")
	white    = lipgloss.Color("#FFFFFF")
)

func main() {
	columns := []table.Column{
		{Title: "Symbol", Width: 8},
		{Title: "Company", Width: 25},
		{Title: "Price", Width: 12},
		{Title: "Change", Width: 22},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithWidth(69),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(deepPink).
		BorderBottom(true).
		Bold(true).
		Foreground(white).
		Background(deepPink).
		Padding(0, 1)

	s.Selected = s.Selected.
		Foreground(deepPink).
		Background(softPink).
		Bold(true)

	s.Cell = s.Cell.
		Foreground(softPink).
		Padding(0, 1)

	t.SetStyles(s)

	m := NewModel()
	m = m.SetTable(t)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
