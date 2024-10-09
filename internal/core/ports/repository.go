package ports

import (
	"github.com/sarthak0714/backend-task-sc/internal/core/domain"
)

type TradeRepository interface {
	AddTrade(trade *domain.Trade) error
	UpdateTrade(id int64, trade *domain.Trade) error
	RemoveTrade(id int64) error
	FetchTrades(userID string) ([]*domain.Trade, error)
}

type PortfolioRepository interface {
	UpsertPortfolio(portfolio *domain.Portfolio) error
	FetchPortfolio(userID string) ([]*domain.Portfolio, error)
}
