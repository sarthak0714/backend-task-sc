package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/sarthak0714/backend-task-sc/internal/adapters/handlers"
	"github.com/sarthak0714/backend-task-sc/internal/adapters/repositories"
	"github.com/sarthak0714/backend-task-sc/internal/config"
	"github.com/sarthak0714/backend-task-sc/internal/core/services"
	"github.com/sarthak0714/backend-task-sc/pkg/utils"
)

func main() {
	// Initialize SQL repository
	cfg := config.LoadConfig()

	tradeRepo, portfolioRepo, err := repositories.NewpgRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error initializing repository: %v", err)
	}

	// Initialize services
	tradeService := services.NewTradeService(tradeRepo, portfolioRepo)
	portfolioService := services.NewPortfolioService(portfolioRepo)

	e := echo.New()
	e.HideBanner = true
	// Middleware
	e.Use(utils.CustomLogger())
	e.Use(middleware.Recover())

	// Initialize handlers
	h := handlers.NewAPIHandler(tradeService, portfolioService)

	e.GET("/", h.Root)

	// Trade Routes
	e.POST("/trades", h.AddTrade)
	e.PUT("/trades/:id", h.UpdateTrade)
	e.DELETE("/trades/:id", h.RemoveTrade)
	e.GET("/trades/:userId", h.FetchTrades)

	//Portfolio Routes
	e.GET("/portfolio/:userId", h.FetchPortfolio)
	e.GET("/returns/:userId", h.FetchReturns)

	e.Logger.Fatal(e.Start(cfg.Port))
}
