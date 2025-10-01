package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/franpfeiffer/t2stock/internal/models"
)

var (
	softPink = lipgloss.Color("#FFB6C1")
	deepPink = lipgloss.Color("#FF1493")
	gray     = lipgloss.Color("#999999")
	red      = lipgloss.Color("#FF6B6B")
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
	case models.InputState:
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
			instructionStyle.Render("Enter your Finnhub API key to get started"),
			"",
			inputStyle.Render(m.textInput.View()),
			"",
			helpStyle.Render("Get your free API key at https://finnhub.io"),
			helpStyle.Render("Free tier: real-time quotes, 60 calls/min"),
			"",
			instructionStyle.Render("Press Enter to continue, Ctrl+C to quit"),
		)
		return containerStyle.Render(content)

	case models.LoadingState:
		loadingStyle := lipgloss.NewStyle().
			Foreground(softPink).
			Align(lipgloss.Center).
			Margin(1, 0)

		dots := strings.Repeat(".", (int(time.Now().UnixMilli()/200))%4)

		content := lipgloss.JoinVertical(lipgloss.Center,
			titleStyle.Render("T2Stock"),
			"",
			loadingStyle.Render("Fetching real-time stock data"+dots),
			loadingStyle.Render("Connecting to Finnhub.io"+dots),
			loadingStyle.Render("Should take just a few seconds."),
		)
		return containerStyle.Render(content)

	case models.StocksState:
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
			helpStyle.Render("Real-time quotes from Finnhub.io"),
		)
		return containerStyle.Render(content)

	case models.ErrorState:
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
