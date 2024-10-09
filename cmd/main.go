package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/sarthak0714/backend-task-sc/docs"
	"github.com/sarthak0714/backend-task-sc/internal/adapters/handlers"
	"github.com/sarthak0714/backend-task-sc/internal/adapters/repositories"
	"github.com/sarthak0714/backend-task-sc/internal/config"
	"github.com/sarthak0714/backend-task-sc/internal/core/services"
	"github.com/sarthak0714/backend-task-sc/pkg/utils"
)

// @title smallcase Backend Task
// @version 1.0
// @description portfolio tracking API.
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

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	e.GET("/status", h.Root)

	// Trade Routes
	e.POST("/trades", h.AddTrade)
	e.PUT("/trades/:id", h.UpdateTrade)
	e.DELETE("/trades/:id", h.RemoveTrade)
	e.GET("/trades/:userId", h.FetchTrades)

	//Portfolio Routes
	e.GET("/portfolio/:userId", h.FetchPortfolio)
	e.GET("/returns", h.FetchReturns)

	// Swagger route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(cfg.Port))
}
