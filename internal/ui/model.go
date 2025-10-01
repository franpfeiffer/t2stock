package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/franpfeiffer/t2stock/internal/models"
	"github.com/franpfeiffer/t2stock/internal/provider"
	"github.com/franpfeiffer/t2stock/internal/storage"
)

type Model struct {
	state     models.AppState
	textInput textinput.Model
	table     table.Model
	apiKey    string
	err       error
	width     int
	height    int
}

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter your Finnhub API key..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	m := Model{
		state:     models.InputState,
		textInput: ti,
	}

	if config, err := storage.LoadConfig(); err == nil && config.APIKey != "" {
		m.apiKey = config.APIKey
		m.state = models.LoadingState
	}

	return m
}

func (m Model) SetTable(t table.Model) Model {
	m.table = t
	return m
}

func (m Model) Init() tea.Cmd {
	if m.state == models.LoadingState {
		return tea.Batch(provider.FetchStockData(m.apiKey), provider.AutoRefresh(m.apiKey, 2*60*time.Second))
	}
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch m.state {
		case models.InputState:
			switch msg.String() {
			case "ctrl+c", "esc":
				return m, tea.Quit
			case "enter":
				if strings.TrimSpace(m.textInput.Value()) != "" {
					m.apiKey = strings.TrimSpace(m.textInput.Value())
					storage.SaveConfig(&models.Config{APIKey: m.apiKey})
					m.state = models.LoadingState
					return m, tea.Batch(provider.FetchStockData(m.apiKey), provider.AutoRefresh(m.apiKey, 2*60*time.Second))
				}
			}

		case models.StocksState:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "r":
				m.state = models.LoadingState
				return m, tea.Batch(provider.FetchStockData(m.apiKey), provider.AutoRefresh(m.apiKey, 2*60*time.Second))
			case "b":
				m.state = models.InputState
				m.textInput.SetValue("")
				m.textInput.Focus()
				return m, textinput.Blink
			case "c":
				storage.SaveConfig(&models.Config{APIKey: ""})
				m.state = models.InputState
				m.textInput.SetValue("")
				m.textInput.Focus()
				return m, textinput.Blink
			}

		case models.ErrorState:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "b":
				m.state = models.InputState
				m.textInput.SetValue("")
				m.textInput.Focus()
				return m, textinput.Blink
			case "r":
				m.state = models.LoadingState
				return m, tea.Batch(provider.FetchStockData(m.apiKey), provider.AutoRefresh(m.apiKey, 2*60*time.Second))
			}
		}

	case models.StockDataMsg:
		m.state = models.StocksState
		rows := make([]table.Row, len(msg))
		for i, stock := range msg {
			changeColor := ""
			if stock.Change != "N/A" && (stock.Change[0] == '+' || stock.Change == "0.00 (0.00%)") {
				changeColor = ""
			} else {
				changeColor = ""
			}
			rows[i] = table.Row{
				stock.Symbol,
				stock.Name,
				stock.Price,
				changeColor + stock.Change,
				stock.Timestamp,
			}
		}
		m.table.SetRows(rows)
		cmd = provider.AutoRefresh(m.apiKey, 2*60*time.Second)

	case models.ErrMsg:
		m.state = models.ErrorState
		m.err = msg
	}

	if m.state == models.InputState {
		m.textInput, cmd = m.textInput.Update(msg)
	} else if m.state == models.StocksState {
		m.table, cmd = m.table.Update(msg)
	}

	return m, cmd
}
