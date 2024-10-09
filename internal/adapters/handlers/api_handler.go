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

func NewAPIHandler(tradeService ports.TradeService, portfolioService ports.PortfolioService) *APIHandler {
	return &APIHandler{tradeService: tradeService, portfolioService: portfolioService}
}

// Root handler
// @Summary Root endpoint
// @Description Returns a simple message indicating the API is working
// @Tags root
// @Produce json
// @Success 200 {object} map[string]string
// @Router / [get]
func (h *APIHandler) Root(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message":   "Works",
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	})
}

// AddTrade adds a new trade
// @Summary Add a new trade
// @Description Adds a new trade to the system
// @Tags trades
// @Accept json
// @Produce json
// @Param trade body domain.Trade true "Trade object"
// @Success 201 {object} domain.Trade
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trades [post]
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

// UpdateTrade updates an existing trade
// @Summary Update a trade
// @Description Updates an existing trade in the system
// @Tags trades
// @Accept json
// @Produce json
// @Param id path int true "Trade ID"
// @Param trade body domain.Trade true "Updated Trade object"
// @Success 200 {object} domain.Trade
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trades/{id} [put]
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

// RemoveTrade deletes an existing trade
// @Summary Remove a trade
// @Description Removes an existing trade from the system
// @Tags trades
// @Produce json
// @Param id path int true "Trade ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trades/{id} [delete]
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

// FetchTrades fetches all trades for a user
// @Summary Fetch user trades
// @Description Fetches all trades for a specific user
// @Tags trades
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {array} domain.Trade
// @Failure 500 {object} map[string]string
// @Router /trades/{userId} [get]
func (h *APIHandler) FetchTrades(c echo.Context) error {
	userID := c.Param("userId")
	trades, err := h.tradeService.FetchTrades(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch trades"})
	}

	return c.JSON(http.StatusOK, trades)
}

// FetchPortfolio fetches portfolio of user
// @Summary Fetch user portfolio
// @Description Fetches the portfolio for a specific user
// @Tags portfolio
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {array} domain.Portfolio
// @Failure 500 {object} map[string]string
// @Router /portfolio/{userId} [get]
func (h *APIHandler) FetchPortfolio(c echo.Context) error {
	userID := c.Param("userId")
	portfolio, err := h.portfolioService.FetchPortfolio(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch portfolio"})
	}

	return c.JSON(http.StatusOK, portfolio)
}

// FetchReturns fetches user returns
// @Summary Fetch user returns
// @Description Fetches the returns for a specific user
// @Tags returns
// @Produce json
// @Param userId query string true "User ID"
// @Success 200 {object} domain.Returns
// @Failure 500 {object} map[string]string
// @Router /returns [get]
func (h *APIHandler) FetchReturns(c echo.Context) error {
	userID := c.QueryParam("userId")
	returns, err := h.portfolioService.FetchReturns(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch returns"})
	}

	return c.JSON(http.StatusOK, returns)
}
