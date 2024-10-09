package domain

import (
	"time"
)

type TradeType string

// Type enum
const (
	Buy  TradeType = "BUY"
	Sell TradeType = "SELL"
)

type Trade struct {
	Id        int64     `json:"id"`
	UserID    string    `gorm:"index" json:"userId"`
	Ticker    string    `json:"ticker"`
	Type      TradeType `json:"type"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

type Portfolio struct {
	UserID          string    `gorm:"index" json:"userId`
	Ticker          string    `json:"ticker"`
	Quantity        int       `json:"quantity"`
	AverageBuyPrice float64   `json:"averageBuyPrice"`
	LastUpdated     time.Time `json:"lastUpdated"`
}

type Returns struct {
	UserID            string  `json:"userId"`
	CumulativeReturns float64 `json:"cumulativeReturns"`
}
