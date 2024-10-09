package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/sarthak0714/backend-task-sc/internal/core/domain"
	"github.com/sarthak0714/backend-task-sc/internal/core/ports"
)

type APIHandler struct {
	tradeService     ports.TradeService
	portfolioService ports.PortfolioService
}

// Returns new Handler service
func NewAPIHandler(tradeService ports.TradeService, portfolioService ports.PortfolioService) *APIHandler {
	return &APIHandler{tradeService: tradeService, portfolioService: portfolioService}
}

// Root handler
func (h *APIHandler) Root(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message":   "Works",
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	})
}

// POST Adds new trade - trade data in body
func (h *APIHandler) AddTrade(c echo.Context) error {
	trade := new(domain.Trade)
	if err := c.Bind(trade); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	trade.Timestamp = time.Now()
	// Basic Validation
	if trade.Quantity <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Quantity must be positive"})
	}
	if trade.Price <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Price must be positive"})
	}
	if trade.Type != domain.Buy && trade.Type != domain.Sell {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid trade type"})
	}

	if err := h.tradeService.AddTrade(trade); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, trade)
}

// PUT Updates exisitng trade - trade data in body
func (h *APIHandler) UpdateTrade(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid trade ID"})
	}

	trade := new(domain.Trade)
	if err := c.Bind(trade); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if err := h.tradeService.UpdateTrade(id, trade); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update trade"})
	}

	return c.JSON(http.StatusOK, trade)
}

// DELETE deletes an existing trade - id in param & data in body
func (h *APIHandler) RemoveTrade(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid trade ID"})
	}

	if err := h.tradeService.RemoveTrade(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to remove trade"})
	}

	return c.NoContent(http.StatusNoContent)
}

// GET fetches all trades or a user - id in param
func (h *APIHandler) FetchTrades(c echo.Context) error {
	userID := c.Param("userId")
	trades, err := h.tradeService.FetchTrades(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch trades"})
	}

	return c.JSON(http.StatusOK, trades)
}

// GET fetches portfolio of user - id in params
func (h *APIHandler) FetchPortfolio(c echo.Context) error {
	userID := c.Param("userId")
	portfolio, err := h.portfolioService.FetchPortfolio(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch portfolio"})
	}

	return c.JSON(http.StatusOK, portfolio)
}

// GET feteches user returns - id in params
func (h *APIHandler) FetchReturns(c echo.Context) error {
	userID := c.QueryParam("userId")
	returns, err := h.portfolioService.FetchReturns(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch returns"})
	}

	return c.JSON(http.StatusOK, returns)
}
