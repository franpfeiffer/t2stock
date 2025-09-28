package main

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	gray = lipgloss.Color("#999999")
	red  = lipgloss.Color("#FF6B6B")
)

func (m Model) View() string {
	containerStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(softPink)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Padding(1, 2).
		Align(lipgloss.Center).
		Foreground(deepPink).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(deepPink)

	switch m.state {
	case InputState:
		inputStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(softPink).
			Padding(1, 2).
			Margin(1, 0)

		helpStyle := lipgloss.NewStyle().
			Foreground(gray).
			Italic(true).
			Align(lipgloss.Center).
			Margin(1, 0)

		instructionStyle := lipgloss.NewStyle().
			Foreground(softPink).
			Align(lipgloss.Center).
			Margin(1, 0)

		content := lipgloss.JoinVertical(lipgloss.Center,
			titleStyle.Render("T2Stock"),
			"",
			instructionStyle.Render("Enter your Polygon.io API key to get started"),
			"",
			inputStyle.Render(m.textInput.View()),
			"",
			helpStyle.Render("Get your free API key at https://polygon.io"),
			helpStyle.Render("Free tier: End-of-day data, 5 calls/min, 2 years history"),
			"",
			instructionStyle.Render("Press Enter to continue, Ctrl+C to quit"),
		)
		return containerStyle.Render(content)

	case LoadingState:
		loadingStyle := lipgloss.NewStyle().
			Foreground(softPink).
			Align(lipgloss.Center).
			Margin(1, 0)

		dots := strings.Repeat(".", (int(time.Now().UnixMilli()/200))%4)

		content := lipgloss.JoinVertical(lipgloss.Center,
			titleStyle.Render("T2Stock"),
			"",
			loadingStyle.Render("Fetching end-of-day stock data"+dots),
			loadingStyle.Render("Connecting to Polygon.io"+dots),
			loadingStyle.Render("Shoud take about a minute."),
		)
		return containerStyle.Render(content)

	case StocksState:
		tableStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(softPink).
			Padding(1, 2).
			Margin(1, 0)

		helpStyle := lipgloss.NewStyle().
			Foreground(gray).
			Italic(true).
			Align(lipgloss.Center).
			Margin(1, 0)

		content := lipgloss.JoinVertical(lipgloss.Center,
			titleStyle.Render("T2Stock"),
			"",
			tableStyle.Render(m.table.View()),
			helpStyle.Render("Press 'r' to refresh, 'b' for new API key, 'c' to clear saved key, 'q' to quit"),
			helpStyle.Render("End-of-day data from Polygon.io"),
		)
		return containerStyle.Render(content)

	case ErrorState:
		errorStyle := lipgloss.NewStyle().
			Foreground(red).
			Align(lipgloss.Center).
			Margin(1, 0)

		helpStyle := lipgloss.NewStyle().
			Foreground(gray).
			Italic(true).
			Align(lipgloss.Center).
			Margin(1, 0)

		content := lipgloss.JoinVertical(lipgloss.Center,
			titleStyle.Render("T2Stock"),
			"",
			errorStyle.Render("Error: "+m.err.Error()),
			"",
			helpStyle.Render("Press 'r' to retry, 'b' for new API key, 'q' to quit"),
		)
		return containerStyle.Render(content)
	}

	return "Unknown state"
}
