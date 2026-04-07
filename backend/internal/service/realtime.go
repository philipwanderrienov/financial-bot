package service

import (
	"fmt"
	"time"

	"finance-agent/backend/internal/client"
	"finance-agent/backend/internal/config"
	"finance-agent/backend/internal/models"
)

type RealtimeService struct {
	finnhub  *client.FinnhubClient
	reco     *RecommendationService
	sectors  []models.Sector
	now      func() time.Time
}

func NewRealtimeService(cfg config.Config, finnhub *client.FinnhubClient, reco *RecommendationService, _ interface{}, sectors []models.Sector) *RealtimeService {
	_ = cfg
	return &RealtimeService{
		finnhub: finnhub,
		reco:    reco,
		sectors: sectors,
		now:     time.Now,
	}
}

func (s *RealtimeService) Start() {}

func (s *RealtimeService) Stop() {}

func (s *RealtimeService) refreshAll() {}

func (s *RealtimeService) refreshSummary() {}

func (s *RealtimeService) refreshWatchlist() {}

func (s *RealtimeService) refreshFilings() {}

func (s *RealtimeService) refreshSectors() {}

func (s *RealtimeService) refreshRecommendations() {}

func (s *RealtimeService) warmUpRecommendations() {}

func firstReason(reasons []string) string { return "" }

func defaultSectors() []models.Sector {
	return []models.Sector{
		{Key: "technology", Label: "Technology", Symbols: []string{"AAPL", "MSFT", "NVDA", "AMD", "INTC", "QCOM"}},
		{Key: "energy", Label: "Energy", Symbols: []string{"XOM", "CVX", "COP", "SLB", "EOG"}},
		{Key: "oil-gas", Label: "Oil & Gas", Symbols: []string{"XOM", "CVX", "COP", "MPC", "VLO", "OXY"}},
	}
}

func (s *RealtimeService) SnapshotForRecommendation(symbol string) models.RecommendationResponse {
	rec, err := s.reco.GetRecommendation(symbol)
	if err == nil {
		return rec
	}

	return models.RecommendationResponse{
		UpdatedAt:  s.now().UTC().Format(time.RFC3339),
		Symbol:     symbol,
		Action:     "hold",
		Confidence: 50,
		Scores: models.RecommendationScores{
			Technical:   50,
			Fundamental: 50,
			News:        50,
			Risk:        50,
		},
		Reasons: []string{"Data sementara tidak tersedia."},
		Sources: models.RecommendationSources{
			MarketData: "Finnhub quote API",
			News:       "Finnhub company news API",
			Filings:    "SEC EDGAR / future integration",
		},
	}
}

func (s *RealtimeService) Seeds() []models.Sector { return s.sectors }

func (s *RealtimeService) Summary() models.SummaryResponse {
	return models.SummaryResponse{UpdatedAt: s.now().UTC().Format(time.RFC3339), Market: models.MarketQuote{Symbol: "AAPL"}, Signals: []models.Signal{}}
}

func (s *RealtimeService) Watchlist() models.WatchlistResponse {
	return models.WatchlistResponse{Items: []models.WatchlistItem{}}
}

func (s *RealtimeService) Filings() models.FilingsResponse {
	return models.FilingsResponse{Items: []models.Filing{}}
}

func (s *RealtimeService) Sectors() []models.Sector { return s.Seeds() }

func (s *RealtimeService) Health() error {
	if s.finnhub == nil {
		return fmt.Errorf("realtime service not initialized")
	}
	return nil
}

func (s *RealtimeService) preferredSymbols() []string { return []string{"AAPL"} }
