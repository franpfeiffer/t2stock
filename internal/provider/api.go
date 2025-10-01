package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/franpfeiffer/t2stock/internal/models"
)

func GetCompanyName(symbol string) string {
	names := map[string]string{
		"AAPL":  "Apple Inc.",
		"GOOGL": "Alphabet Inc.",
		"AMZN":  "Amazon.com Inc.",
		"SONY":  "Sony Group Corp.",
		"SPY":   "SPDR S&P 500 ETF",
	}
	if name, exists := names[symbol]; exists {
		return name
	}
	return symbol
}

func FetchFinnhubStock(symbol, apiKey string) (models.Stock, error) {
	url := fmt.Sprintf("https://finnhub.io/api/v1/quote?symbol=%s&token=%s", symbol, apiKey)
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.Stock{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "T2Stock/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return models.Stock{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 401 {
		return models.Stock{}, fmt.Errorf("invalid API key")
	}
	if resp.StatusCode == 429 {
		return models.Stock{}, fmt.Errorf("rate limit exceeded - free tier allows 60 calls/min")
	}
	if resp.StatusCode != 200 {
		return models.Stock{}, fmt.Errorf("API error %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Stock{}, err
	}
	var quote struct {
		Current       float64 `json:"c"`
		Change        float64 `json:"d"`
		ChangePercent float64 `json:"dp"`
		Timestamp     int64   `json:"t"`
	}
	if err := json.Unmarshal(body, &quote); err != nil {
		return models.Stock{}, err
	}
	price := fmt.Sprintf("$%.2f", quote.Current)
	changeStr := fmt.Sprintf("%.2f (%.2f%%)", quote.Change, quote.ChangePercent)
	if quote.Change >= 0 {
		changeStr = "+" + changeStr
	}
	ts := time.Unix(quote.Timestamp, 0).In(time.Local).Format("2006-01-02 15:04:05")
	return models.Stock{
		Symbol:    symbol,
		Name:      GetCompanyName(symbol),
		Price:     price,
		Change:    changeStr,
		Timestamp: ts,
	}, nil
}

func FetchStockData(apiKey string) tea.Cmd {
	return func() tea.Msg {
		symbols := []string{"AAPL", "GOOGL", "AMZN", "SONY", "SPY"}
		var stockData []models.Stock
		for _, symbol := range symbols {
			stock, err := FetchFinnhubStock(symbol, apiKey)
			if err != nil {
				if len(stockData) == 0 {
					return models.ErrMsg(err)
				}
				stock = models.Stock{
					Symbol:    symbol,
					Name:      GetCompanyName(symbol),
					Price:     "N/A",
					Change:    "N/A",
					Timestamp: "N/A",
				}
			}
			stockData = append(stockData, stock)
			time.Sleep(1 * time.Second)
		}
		return models.StockDataMsg(stockData)
	}
}

func AutoRefresh(apiKey string, interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(time.Time) tea.Msg {
		return FetchStockData(apiKey)()
	})
}
