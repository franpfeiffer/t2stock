package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func GetLastTradingDay() string {
	now := time.Now()

	for i := 1; i <= 7; i++ {
		date := now.AddDate(0, 0, -i)
		weekday := date.Weekday()

		if weekday != time.Saturday && weekday != time.Sunday {
			return date.Format("2006-01-02")
		}
	}

	return now.AddDate(0, 0, -3).Format("2006-01-02")
}

func FetchPolygonStock(symbol, apiKey string) (Stock, error) {
	lastTradingDay := GetLastTradingDay()

	url := fmt.Sprintf("https://api.polygon.io/v2/aggs/ticker/%s/range/1/day/%s/%s?adjusted=true&sort=desc&limit=1&apikey=%s",
		symbol, lastTradingDay, lastTradingDay, apiKey)

	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Stock{}, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "T2Stock/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return Stock{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return Stock{}, fmt.Errorf("invalid API key")
	}
	if resp.StatusCode == 403 {
		return Stock{}, fmt.Errorf("access denied - upgrade plan needed")
	}
	if resp.StatusCode == 429 {
		return Stock{}, fmt.Errorf("rate limit exceeded - free tier allows 5 calls/min")
	}

	if resp.StatusCode != 200 {
		return Stock{}, fmt.Errorf("API error %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Stock{}, err
	}

	var aggResp struct {
		Status  string `json:"status"`
		Results []struct {
			Open   float64 `json:"o"`
			High   float64 `json:"h"`
			Low    float64 `json:"l"`
			Close  float64 `json:"c"`
			Volume float64 `json:"v"`
		} `json:"results"`
	}

	if err := json.Unmarshal(body, &aggResp); err != nil {
		return Stock{}, err
	}

	if aggResp.Status != "OK" || len(aggResp.Results) == 0 {
		return Stock{}, fmt.Errorf("no data available")
	}

	result := aggResp.Results[0]
	currentPrice := result.Close

	change := currentPrice - result.Open
	changePercent := (change / result.Open) * 100

	price := fmt.Sprintf("$%.2f", currentPrice)
	changeStr := fmt.Sprintf("%.2f (%.2f%%)", change, changePercent)
	if change >= 0 {
		changeStr = "+" + changeStr
	}

	return Stock{
		Symbol: symbol,
		Name:   GetCompanyName(symbol),
		Price:  price,
		Change: changeStr,
	}, nil
}

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

func FetchStockData(apiKey string) tea.Cmd {
	return func() tea.Msg {
		symbols := []string{
			"AAPL",
			"GOOGL",
			"AMZN",
			"SONY",
			"SPY",
		}

		var stockData []Stock

		for _, symbol := range symbols {
			stock, err := FetchPolygonStock(symbol, apiKey)

			if err != nil {
				if len(stockData) == 0 {
					return ErrMsg(err)
				}
				stock = Stock{
					Symbol: symbol,
					Name:   GetCompanyName(symbol),
					Price:  "N/A",
					Change: "N/A",
				}
			}

			stockData = append(stockData, stock)
			time.Sleep(12 * time.Second)
		}

		return StockDataMsg(stockData)
	}
}
