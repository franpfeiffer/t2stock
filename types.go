package main

type AppState int

const (
	InputState AppState = iota
	LoadingState
	StocksState
	ErrorState
)

type Config struct {
	APIKey string `json:"api_key"`
}

type Stock struct {
	Symbol string
	Price  string
	Change string
	Name   string
}

type PolygonResponse struct {
	Status     string `json:"status"`
	RequestID  string `json:"request_id"`
	Count      int    `json:"count"`
	Results    []struct {
		Value            float64 `json:"value"`
		Last             struct {
			Conditions      []int   `json:"conditions"`
			Exchange        int     `json:"exchange"`
			Price           float64 `json:"price"`
			Sip_timestamp   int64   `json:"sip_timestamp"`
			Timeframe       string  `json:"timeframe"`
		} `json:"last_quote"`
		LastTrade struct {
			Conditions      []int   `json:"conditions"`
			Exchange        int     `json:"exchange"`
			Price           float64 `json:"price"`
			Sip_timestamp   int64   `json:"sip_timestamp"`
			Timeframe       string  `json:"timeframe"`
			Size            int     `json:"size"`
		} `json:"last_trade"`
		Market_status   string  `json:"market_status"`
		Name            string  `json:"name"`
		Type            string  `json:"type"`
		Session         struct {
			Change          float64 `json:"change"`
			Change_percent  float64 `json:"change_percent"`
			Early_trading_change float64 `json:"early_trading_change"`
			Early_trading_change_percent float64 `json:"early_trading_change_percent"`
			Close           float64 `json:"close"`
			High            float64 `json:"high"`
			Low             float64 `json:"low"`
			Open            float64 `json:"open"`
			Previous_close  float64 `json:"previous_close"`
		} `json:"session"`
		Ticker          string  `json:"ticker"`
		Updated         int64   `json:"updated"`
	} `json:"results"`
}

type StockDataMsg []Stock
type ErrMsg error
