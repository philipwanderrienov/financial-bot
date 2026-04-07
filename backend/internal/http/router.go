package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"finance-agent/backend/internal/client"
	"finance-agent/backend/internal/config"
	"finance-agent/backend/internal/models"
	"finance-agent/backend/internal/service"
)

func SetupRouter(cfg config.Config) *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5174")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	finnhubClient := client.NewFinnhubClient(cfg.FinnhubAPIKey)
	recoService := service.NewRecommendationService(finnhubClient)
	realtimeService := service.NewRealtimeService(cfg, finnhubClient, recoService, nil, defaultSectors())
	SetRealtimeService(realtimeService)
	realtimeService.Start()

	router.GET("/", SwaggerUI)
	router.GET("/swagger", SwaggerUI)
	router.GET("/swagger.json", SwaggerJSON)

	api := router.Group("/api")
	{
		api.GET("/health", Health)
		api.GET("/summary", Summary)
		api.GET("/watchlist", Watchlist)
		api.GET("/filings", Filings)
		api.GET("/recommendation", Recommendation)
	}

	return router
}

func defaultSectors() []models.Sector {
	return []models.Sector{
		{Key: "technology", Label: "Technology", Symbols: []string{"AAPL", "MSFT", "NVDA", "AMD", "INTC", "QCOM"}},
		{Key: "energy", Label: "Energy", Symbols: []string{"XOM", "CVX", "COP", "SLB", "EOG"}},
		{Key: "oil-gas", Label: "Oil & Gas", Symbols: []string{"XOM", "CVX", "COP", "MPC", "VLO", "OXY"}},
	}
}
