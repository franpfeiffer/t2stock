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
	Symbol    string
	Name      string
	Price     string
	Change    string
	Timestamp string
}

type StockDataMsg []Stock
type ErrMsg error
