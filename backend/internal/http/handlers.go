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

	if realtimeService == nil {
		c.JSON(http.StatusOK, models.RecommendationResponse{
			UpdatedAt:  "2026-04-06T00:00:00Z",
			Symbol:     symbol,
			Action:     "hold",
			Confidence: 50,
			Scores: models.RecommendationScores{
				Technical:   50,
				Fundamental: 50,
				News:        50,
				Risk:        50,
			},
			Reasons: []string{"Data sementara belum tersedia."},
			Sources: models.RecommendationSources{
				MarketData: "Finnhub quote API",
				News:       "Finnhub company news API",
				Filings:    "SEC EDGAR / future integration",
			},
		})
		return
	}

	c.JSON(http.StatusOK, realtimeService.SnapshotForRecommendation(symbol))
}