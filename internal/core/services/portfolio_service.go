package services

import (
	"fmt"

	"github.com/sarthak0714/backend-task-sc/internal/core/domain"
	"github.com/sarthak0714/backend-task-sc/internal/core/ports"
)

type portfolioService struct {
	portfolioRepo ports.PortfolioRepository
}

// Creates a new Portfolio Service
func NewPortfolioService(portfolioRepo ports.PortfolioRepository) ports.PortfolioService {
	return &portfolioService{portfolioRepo: portfolioRepo}
}

// Fetches a user portfolio
func (s *portfolioService) FetchPortfolio(userID string) ([]*domain.Portfolio, error) {
	return s.portfolioRepo.FetchPortfolio(userID)
}

// Fetches users returns (Based on given Logic)
func (s *portfolioService) FetchReturns(userID string) (*domain.Returns, error) {
	portfolio, err := s.portfolioRepo.FetchPortfolio(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch portfolio: %w", err)
	}

	if len(portfolio) == 0 {
		return &domain.Returns{UserID: userID, CumulativeReturns: 0}, nil
	}

	var cumulativeReturns float64
	for _, security := range portfolio {
		currentPrice := 100.0 // Assuming current price is always 100
		returns := (currentPrice - security.AverageBuyPrice) * float64(security.Quantity)
		cumulativeReturns += returns
	}

	return &domain.Returns{UserID: userID, CumulativeReturns: cumulativeReturns}, nil
}
