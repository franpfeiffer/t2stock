package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	state     AppState
	textInput textinput.Model
	table     table.Model
	apiKey    string
	err       error
	width     int
	height    int
}

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter your Polygon.io API key..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	m := Model{
		state:     InputState,
		textInput: ti,
	}

	if config, err := LoadConfig(); err == nil && config.APIKey != "" {
		m.apiKey = config.APIKey
		m.state = LoadingState
	}

	return m
}

func (m Model) SetTable(t table.Model) Model {
	m.table = t
	return m
}

func (m Model) Init() tea.Cmd {
	if m.state == LoadingState {
		return FetchStockData(m.apiKey)
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
		case InputState:
			switch msg.String() {
			case "ctrl+c", "esc":
				return m, tea.Quit
			case "enter":
				if strings.TrimSpace(m.textInput.Value()) != "" {
					m.apiKey = strings.TrimSpace(m.textInput.Value())

					config := &Config{APIKey: m.apiKey}
					SaveConfig(config)

					m.state = LoadingState
					return m, FetchStockData(m.apiKey)
				}
			}

		case StocksState:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "r":
				m.state = LoadingState
				return m, FetchStockData(m.apiKey)
			case "b":
				m.state = InputState
				m.textInput.SetValue("")
				m.textInput.Focus()
				return m, textinput.Blink
			case "c":
				config := &Config{APIKey: ""}
				SaveConfig(config)
				m.state = InputState
				m.textInput.SetValue("")
				m.textInput.Focus()
				return m, textinput.Blink
			}

		case ErrorState:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "b":
				m.state = InputState
				m.textInput.SetValue("")
				m.textInput.Focus()
				return m, textinput.Blink
			case "r":
				m.state = LoadingState
				return m, FetchStockData(m.apiKey)
			}
		}

	case StockDataMsg:
		m.state = StocksState
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
			}
		}
		m.table.SetRows(rows)

	case ErrMsg:
		m.state = ErrorState
		m.err = msg
	}

	if m.state == InputState {
		m.textInput, cmd = m.textInput.Update(msg)
	} else if m.state == StocksState {
		m.table, cmd = m.table.Update(msg)
	}

	return m, cmd
}
