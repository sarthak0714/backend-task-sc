package repositories

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/sarthak0714/backend-task-sc/internal/core/domain"
	"github.com/sarthak0714/backend-task-sc/internal/core/ports"
)

type pgRepository struct {
	db *gorm.DB
}

// Creates and initializes new Repositories
func NewpgRepository(dbUrl string) (ports.TradeRepository, ports.PortfolioRepository, error) {
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	repo := &pgRepository{db: db}

	db.AutoMigrate(&domain.Trade{}, &domain.Portfolio{})
	// at the tables are in same db but created isolated repos for scalablity
	return repo, repo, nil
}

// Adds a new Trade (with all validations)
func (r *pgRepository) AddTrade(trade *domain.Trade) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Fetch current portfolio item
		var portfolio domain.Portfolio
		if err := tx.Where("user_id = ? AND ticker = ?", trade.UserID, trade.Ticker).First(&portfolio).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				portfolio = domain.Portfolio{
					UserID:          trade.UserID,
					Ticker:          trade.Ticker,
					Quantity:        0,
					AverageBuyPrice: 0,
				}
				er := tx.Create(portfolio).Error
				if er != nil {
					return er
				}
			} else {
				return err
			}
		}

		// Update portfolio based on trade type
		switch trade.Type {
		case domain.Buy:
			newQuantity := portfolio.Quantity + trade.Quantity
			newTotalValue := (portfolio.AverageBuyPrice * float64(portfolio.Quantity)) + (trade.Price * float64(trade.Quantity))
			portfolio.Quantity = newQuantity
			if newQuantity > 0 {
				portfolio.AverageBuyPrice = newTotalValue / float64(newQuantity)
			} else {
				portfolio.AverageBuyPrice = 0
			}
		case domain.Sell:
			if portfolio.Quantity < trade.Quantity {
				return fmt.Errorf("insufficient quantity for sell trade")
			}
			portfolio.Quantity -= trade.Quantity
			// no change to AverageBuyPrice when selling
		}

		portfolio.LastUpdated = trade.Timestamp

		// Save or update portfolio
		if err := tx.Where("user_id = ? AND ticker = ?", trade.UserID, trade.Ticker).Save(&portfolio).Error; err != nil {
			return err
		}

		// Add trade
		return tx.Create(trade).Error
	})
}

// Updates a existing Trade (with all validations)
func (r *pgRepository) UpdateTrade(id int64, updatedTrade *domain.Trade) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Fetch original trade
		var originalTrade domain.Trade
		if err := tx.First(&originalTrade, id).Error; err != nil {
			return err
		}

		// Fetch current portfolio item
		var portfolio domain.Portfolio
		if err := tx.Where("user_id = ? AND ticker = ?", originalTrade.UserID, originalTrade.Ticker).First(&portfolio).Error; err != nil {
			return err
		}

		// Revert the effect of the original trade
		switch originalTrade.Type {
		case domain.Buy:
			portfolio.Quantity -= originalTrade.Quantity
			if portfolio.Quantity > 0 {
				portfolio.AverageBuyPrice = (portfolio.AverageBuyPrice*float64(portfolio.Quantity+originalTrade.Quantity) - float64(originalTrade.Quantity)*originalTrade.Price) / float64(portfolio.Quantity)
			} else {
				portfolio.AverageBuyPrice = 0
			}
		case domain.Sell:
			portfolio.Quantity += originalTrade.Quantity
		}

		// Apply the updated trade
		switch updatedTrade.Type {
		case domain.Buy:
			newQuantity := portfolio.Quantity + updatedTrade.Quantity
			newTotalValue := (portfolio.AverageBuyPrice * float64(portfolio.Quantity)) + (updatedTrade.Price * float64(updatedTrade.Quantity))
			portfolio.Quantity = newQuantity
			if newQuantity > 0 {
				portfolio.AverageBuyPrice = newTotalValue / float64(newQuantity)
			} else {
				portfolio.AverageBuyPrice = 0
			}
		case domain.Sell:
			if portfolio.Quantity < updatedTrade.Quantity {
				return fmt.Errorf("insufficient quantity for updated sell trade")
			}
			portfolio.Quantity -= updatedTrade.Quantity
		}

		portfolio.LastUpdated = updatedTrade.Timestamp

		// Update portfolio
		if err := tx.Where("user_id = ? AND ticker = ?", originalTrade.UserID, originalTrade.Ticker).Save(&portfolio).Error; err != nil {
			return err
		}

		// Update trade
		return tx.Model(&domain.Trade{}).Where("id = ?", updatedTrade.Id).Updates(updatedTrade).Error
	})
}

// Removes a Trade (with all validations)
func (r *pgRepository) RemoveTrade(id int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Fetch the trade to be removed
		var trade domain.Trade
		if err := tx.First(&trade, id).Error; err != nil {
			return err
		}

		// Fetch current portfolio item
		var portfolio domain.Portfolio
		if err := tx.Where("user_id = ? AND ticker = ?", trade.UserID, trade.Ticker).First(&portfolio).Error; err != nil {
			return err
		}

		// Revert the effect of the trade
		switch trade.Type {
		case domain.Buy:
			portfolio.Quantity -= trade.Quantity
			if portfolio.Quantity > 0 {
				portfolio.AverageBuyPrice = (portfolio.AverageBuyPrice*float64(portfolio.Quantity+trade.Quantity) - float64(trade.Quantity)*trade.Price) / float64(portfolio.Quantity)
			} else {
				portfolio.AverageBuyPrice = 0
			}
		case domain.Sell:
			portfolio.Quantity += trade.Quantity
		}

		if portfolio.Quantity < 0 {
			return fmt.Errorf("removing this trade would result in negative quantity")
		}

		portfolio.LastUpdated = time.Now()

		if err := tx.Where("user_id = ? AND ticker = ?", trade.UserID, trade.Ticker).Save(&portfolio).Error; err != nil {
			return err
		}

		// Remove trade
		return tx.Delete(&trade).Error
	})
}

// Fetch all trades for a user
func (r *pgRepository) FetchTrades(userID string) ([]*domain.Trade, error) {
	var trades []*domain.Trade
	err := r.db.Where("user_id = ?", userID).Find(&trades).Error
	return trades, err
}

// Fetch Portfolio
func (r *pgRepository) FetchPortfolio(userID string) ([]*domain.Portfolio, error) {
	var portfolio []*domain.Portfolio
	err := r.db.Where("user_id = ?", userID).Find(&portfolio).Error
	return portfolio, err
}
