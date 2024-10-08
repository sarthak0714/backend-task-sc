package ports

import (
	"context"

	"github.com/sarthak0714/backend-task-sc/internal/core/domain"
)

type TradeRepository interface {
	AddTrade(ctx context.Context, trade *domain.Trade) error
	UpdateTrade(ctx context.Context, id int64, trade *domain.Trade) error
	RemoveTrade(ctx context.Context, id int64) error
	FetchTrades(ctx context.Context, userID string) ([]*domain.Trade, error)
}

type PortfolioRepository interface {
	UpsertPortfolio(ctx context.Context, portfolio *domain.Portfolio) error
	FetchPortfolio(ctx context.Context, userID string) ([]*domain.Portfolio, error)
}
