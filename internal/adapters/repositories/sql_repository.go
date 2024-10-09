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

func NewpgRepository(dbUrl string) (ports.TradeRepository, ports.PortfolioRepository, error) {
	db, err := gorm.Open(postgres.Open(dbUrl))
	if err != nil {
		return nil, nil, err
	}

	repo := &pgRepository{db: db}

	db.AutoMigrate(&domain.Trade{}, &domain.Portfolio{})

	// Add indexes
	// db.Model(&domain.Trade{}).AddIndex("idx_user_id_ticker", "user_id", "ticker")
	// db.Model(&domain.Portfolio{}).AddIndex("idx_user_id_ticker", "user_id", "ticker")

	return repo, repo, nil
}

func (r *pgRepository) AddTrade(trade *domain.Trade) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Fetch current portfolio item for the specific security
		var portfolio domain.Portfolio
		if err := tx.Where("user_id = ? AND ticker = ?", trade.UserID, trade.Ticker).First(&portfolio).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
			// If portfolio item doesn't exist, create a new one
			portfolio = domain.Portfolio{UserID: trade.UserID, Ticker: trade.Ticker, Quantity: 0, AverageBuyPrice: 0}
		}

		// Ensure that the portfolio is updated based on the correct ticker
		// Calculate new quantity and average buy price
		newQuantity := portfolio.Quantity
		newTotalValue := float64(portfolio.Quantity) * portfolio.AverageBuyPrice

		if trade.Type == domain.Buy {
			newQuantity += trade.Quantity
			newTotalValue += float64(trade.Quantity) * trade.Price
		} else if trade.Type == domain.Sell {
			// Validate that quantity doesn't go negative
			if portfolio.Quantity < trade.Quantity {
				return fmt.Errorf("insufficient quantity for sell trade")
			}
			newQuantity -= trade.Quantity
		}

		// Update portfolio
		portfolio.Quantity = newQuantity
		if newQuantity > 0 {
			portfolio.AverageBuyPrice = newTotalValue / float64(newQuantity)
		} else {
			portfolio.AverageBuyPrice = 0
		}
		portfolio.LastUpdated = trade.Timestamp

		if err := tx.Save(&portfolio).Error; err != nil {
			return err
		}

		// Add trade
		return tx.Create(trade).Error
	})
}

func (r *pgRepository) UpdateTrade(id int64, updatedTrade *domain.Trade) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Fetch original trade
		var originalTrade domain.Trade
		if err := tx.First(&originalTrade, id).Error; err != nil {
			return err
		}

		// Fetch current portfolio item for the specific security
		var portfolio domain.Portfolio
		if err := tx.Where("user_id = ? AND ticker = ?", originalTrade.UserID, originalTrade.Ticker).First(&portfolio).Error; err != nil {
			return err
		}

		// Revert the effect of the original trade
		if originalTrade.Type == domain.Buy {
			portfolio.Quantity -= originalTrade.Quantity
			portfolio.AverageBuyPrice = (portfolio.AverageBuyPrice*float64(portfolio.Quantity) - float64(originalTrade.Quantity)*originalTrade.Price) / float64(portfolio.Quantity-originalTrade.Quantity)
		} else if originalTrade.Type == domain.Sell {
			portfolio.Quantity += originalTrade.Quantity
		}

		// Apply the updated trade
		if updatedTrade.Type == domain.Buy {
			portfolio.Quantity += updatedTrade.Quantity
			portfolio.AverageBuyPrice = (portfolio.AverageBuyPrice*float64(portfolio.Quantity) + float64(updatedTrade.Quantity)*updatedTrade.Price) / float64(portfolio.Quantity+updatedTrade.Quantity)
		} else if updatedTrade.Type == domain.Sell {
			portfolio.Quantity -= updatedTrade.Quantity
			// Validate that quantity doesn't go negative
			if portfolio.Quantity < 0 {
				return fmt.Errorf("insufficient quantity for updated sell trade")
			}
		}

		portfolio.LastUpdated = updatedTrade.Timestamp

		// Update portfolio
		if err := tx.Save(&portfolio).Error; err != nil {
			return err
		}

		// Update trade
		return tx.Model(&domain.Trade{}).Where("id = ?", id).Updates(updatedTrade).Error
	})
}

func (r *pgRepository) RemoveTrade(id int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Fetch the trade to be removed
		var trade domain.Trade
		if err := tx.First(&trade, id).Error; err != nil {
			return err
		}

		// Fetch current portfolio item for the specific security
		var portfolio domain.Portfolio
		if err := tx.Where("user_id = ? AND ticker = ?", trade.UserID, trade.Ticker).First(&portfolio).Error; err != nil {
			return err
		}

		// Revert the effect of the trade
		if trade.Type == domain.Buy {
			portfolio.Quantity -= trade.Quantity
			if portfolio.Quantity > 0 {
				portfolio.AverageBuyPrice = (portfolio.AverageBuyPrice*float64(portfolio.Quantity+trade.Quantity) - float64(trade.Quantity)*trade.Price) / float64(portfolio.Quantity)
			} else {
				portfolio.AverageBuyPrice = 0
			}
		} else if trade.Type == domain.Sell {
			portfolio.Quantity += trade.Quantity
			// Recalculate average buy price for sell trades
			portfolio.AverageBuyPrice = (portfolio.AverageBuyPrice*float64(portfolio.Quantity-trade.Quantity) + float64(trade.Quantity)*trade.Price) / float64(portfolio.Quantity)
		}

		// Validate that quantity doesn't go negative
		if portfolio.Quantity < 0 {
			return fmt.Errorf("removing this trade would result in negative quantity")
		}

		portfolio.LastUpdated = time.Now()

		// Update portfolio
		if err := tx.Save(&portfolio).Error; err != nil {
			return err
		}

		// Remove trade
		return tx.Delete(&trade).Error
	})
}

func (r *pgRepository) FetchTrades(userID string) ([]*domain.Trade, error) {
	var trades []*domain.Trade
	err := r.db.Where("user_id = ?", userID).Find(&trades).Error
	return trades, err
}

func (r *pgRepository) UpsertPortfolio(portfolio *domain.Portfolio) error {
	return r.db.Save(portfolio).Error
}

func (r *pgRepository) FetchPortfolio(userID string) ([]*domain.Portfolio, error) {
	var portfolio []*domain.Portfolio
	err := r.db.Where("user_id = ?", userID).Find(&portfolio).Error
	return portfolio, err
}
