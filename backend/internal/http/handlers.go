package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"finance-agent/backend/internal/models"
	"finance-agent/backend/internal/service"
)

var realtimeService *service.RealtimeService

func SetRealtimeService(s *service.RealtimeService) {
	realtimeService = s
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func Summary(c *gin.Context) {
	if realtimeService == nil {
		c.JSON(http.StatusOK, models.SummaryResponse{
			UpdatedAt: "2026-04-06T00:00:00Z",
			Market:    models.MarketQuote{Symbol: "AAPL"},
			Signals:   []models.Signal{},
		})
		return
	}
	c.JSON(http.StatusOK, realtimeService.Summary())
}

func Watchlist(c *gin.Context) {
	if realtimeService == nil {
		c.JSON(http.StatusOK, models.WatchlistResponse{Items: []models.WatchlistItem{}})
		return
	}
	c.JSON(http.StatusOK, realtimeService.Watchlist())
}

func Filings(c *gin.Context) {
	if realtimeService == nil {
		c.JSON(http.StatusOK, models.FilingsResponse{Items: []models.Filing{}})
		return
	}
	c.JSON(http.StatusOK, realtimeService.Filings())
}

func Recommendation(c *gin.Context) {
	symbol := c.Query("symbol")
	if symbol == "" {
		symbol = "AAPL"
	}

	c.JSON(http.StatusOK, realtimeService.SnapshotForRecommendation(symbol))
}
