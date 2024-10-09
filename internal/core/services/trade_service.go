package services

import (
	"github.com/sarthak0714/backend-task-sc/internal/core/domain"
	"github.com/sarthak0714/backend-task-sc/internal/core/ports"
)

type tradeService struct {
	tradeRepo     ports.TradeRepository
	portfolioRepo ports.PortfolioRepository
}

// Creates a new Trade Service
func NewTradeService(tradeRepo ports.TradeRepository, portfolioRepo ports.PortfolioRepository) ports.TradeService {
	return &tradeService{tradeRepo: tradeRepo, portfolioRepo: portfolioRepo}
}

// Adds new Trade
func (s *tradeService) AddTrade(trade *domain.Trade) error {
	return s.tradeRepo.AddTrade(trade)

}

// Updates a existing trade
func (s *tradeService) UpdateTrade(id int64, trade *domain.Trade) error {
	return s.tradeRepo.UpdateTrade(id, trade)
}

// Removes a trade based on ID
func (s *tradeService) RemoveTrade(id int64) error {
	return s.tradeRepo.RemoveTrade(id)
}

// Fetches all trades for a User
func (s *tradeService) FetchTrades(userID string) ([]*domain.Trade, error) {
	return s.tradeRepo.FetchTrades(userID)
}
