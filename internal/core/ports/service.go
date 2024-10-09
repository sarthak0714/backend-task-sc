package ports

import (
	"github.com/sarthak0714/backend-task-sc/internal/core/domain"
)

type TradeService interface {
	AddTrade(trade *domain.Trade) error
	UpdateTrade(id int64, trade *domain.Trade) error
	RemoveTrade(id int64) error
	FetchTrades(userID string) ([]*domain.Trade, error)
}

type PortfolioService interface {
	FetchPortfolio(userID string) ([]*domain.Portfolio, error)
	FetchReturns(userID string) (*domain.Returns, error)
}
