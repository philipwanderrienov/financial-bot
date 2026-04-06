package service

import (
	"fmt"
	"time"

	"finance-agent/backend/internal/cache"
	"finance-agent/backend/internal/client"
	"finance-agent/backend/internal/config"
	"finance-agent/backend/internal/models"
)

type RealtimeService struct {
	finnhub   *client.FinnhubClient
	reco      *RecommendationService
	store     *cache.SnapshotStore
	sectors   []models.Sector
	interval  time.Duration
	now       func() time.Time
	stopCh    chan struct{}
	stoppedCh chan struct{}
}

func NewRealtimeService(cfg config.Config, finnhub *client.FinnhubClient, reco *RecommendationService, store *cache.SnapshotStore, sectors []models.Sector) *RealtimeService {
	interval := 15 * time.Second
	if cfg.RefreshSeconds > 0 {
		interval = time.Duration(cfg.RefreshSeconds) * time.Second
	}

	return &RealtimeService{
		finnhub:   finnhub,
		reco:      reco,
		store:     store,
		sectors:   sectors,
		interval:  interval,
		now:       time.Now,
		stopCh:    make(chan struct{}),
		stoppedCh: make(chan struct{}),
	}
}

func (s *RealtimeService) Start() {
	s.refreshAll()

	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()
		defer close(s.stoppedCh)

		for {
			select {
			case <-ticker.C:
				s.refreshAll()
			case <-s.stopCh:
				return
			}
		}
	}()
}

func (s *RealtimeService) Stop() {
	select {
	case <-s.stopCh:
		return
	default:
		close(s.stopCh)
		<-s.stoppedCh
	}
}

func (s *RealtimeService) refreshAll() {
	s.refreshSummary()
	s.refreshWatchlist()
	s.refreshFilings()
	s.refreshSectors()
	s.refreshRecommendations()
}

func (s *RealtimeService) refreshSummary() {
	symbols := []string{"AAPL", "MSFT", "NVDA"}
	items := make([]models.Signal, 0, len(symbols))
	market := models.MarketQuote{Symbol: "AAPL"}

	for i, symbol := range symbols {
		rec, err := s.reco.GetRecommendation(symbol)
		if err != nil {
			if i == 0 && s.store.HasSummary() {
				return
			}
			continue
		}
		s.store.UpdateRecommendation(symbol, rec)
		items = append(items, models.Signal{
			Symbol:     symbol,
			Signal:     rec.Action,
			Confidence: rec.Confidence / 100,
			Reason:     firstReason(rec.Reasons),
		})
		if symbol == "AAPL" {
			market = models.MarketQuote{Symbol: symbol}
		}
	}

	if len(items) == 0 && s.store.HasSummary() {
		return
	}

	summary := models.SummaryResponse{
		UpdatedAt: s.now().UTC().Format(time.RFC3339),
		Market:    market,
		Signals:   items,
	}
	s.store.UpdateSummary(summary)
}

func (s *RealtimeService) refreshWatchlist() {
	items := []models.WatchlistItem{}
	for _, symbol := range []string{"AAPL", "MSFT", "NVDA"} {
		rec, err := s.reco.GetRecommendation(symbol)
		if err != nil {
			continue
		}
		items = append(items, models.WatchlistItem{
			Symbol:        symbol,
			Name:          symbol,
			Price:         rec.Confidence,
			ChangePercent: rec.Scores.Technical - 50,
			Signal:        rec.Action,
		})
	}
	if len(items) == 0 && s.store.HasWatchlist() {
		return
	}
	s.store.UpdateWatchlist(models.WatchlistResponse{Items: items})
}

func (s *RealtimeService) refreshFilings() {
	filings := models.FilingsResponse{
		Items: []models.Filing{
			{
				Symbol:      "AAPL",
				Title:       "Latest quarterly filing available",
				Source:      "SEC",
				PublishedAt: s.now().UTC().Format(time.RFC3339),
				URL:         "https://www.sec.gov",
			},
		},
	}
	if s.store.HasFilings() {
		s.store.UpdateFilings(filings)
		return
	}
	s.store.UpdateFilings(filings)
}

func (s *RealtimeService) refreshSectors() {
	if len(s.sectors) == 0 {
		return
	}
	s.store.UpdateSectors(s.sectors)
}

func (s *RealtimeService) refreshRecommendations() {
	for _, sector := range s.sectors {
		for _, symbol := range sector.Symbols {
			rec, err := s.reco.GetRecommendation(symbol)
			if err != nil {
				continue
			}
			s.store.UpdateRecommendation(symbol, rec)
		}
	}
}

func firstReason(reasons []string) string {
	if len(reasons) == 0 {
		return ""
	}
	return reasons[0]
}

func defaultSectors() []models.Sector {
	return []models.Sector{
		{Key: "technology", Label: "Technology", Symbols: []string{"AAPL", "MSFT", "NVDA", "AMD", "INTC", "QCOM"}},
		{Key: "energy", Label: "Energy", Symbols: []string{"XOM", "CVX", "COP", "SLB", "EOG"}},
		{Key: "oil-gas", Label: "Oil & Gas", Symbols: []string{"XOM", "CVX", "COP", "MPC", "VLO", "OXY"}},
	}
}

func (s *RealtimeService) SnapshotForRecommendation(symbol string) models.RecommendationResponse {
	if rec, ok := s.store.Recommendation(symbol); ok {
		return rec
	}
	rec, err := s.reco.GetRecommendation(symbol)
	if err == nil {
		s.store.UpdateRecommendation(symbol, rec)
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
		Reasons: []string{"Data sementara belum tersedia, menunggu refresh berikutnya."},
		Sources: models.RecommendationSources{
			MarketData: "Finnhub quote API",
			News:       "Finnhub company news API",
			Filings:    "SEC EDGAR / future integration",
		},
	}
}

func (s *RealtimeService) Seeds() []models.Sector {
	if len(s.sectors) > 0 {
		return s.sectors
	}
	return defaultSectors()
}

func (s *RealtimeService) Summary() models.SummaryResponse {
	if s.store.HasSummary() {
		return s.store.Summary()
	}
	return models.SummaryResponse{UpdatedAt: s.now().UTC().Format(time.RFC3339), Market: models.MarketQuote{Symbol: "AAPL"}, Signals: []models.Signal{}}
}

func (s *RealtimeService) Watchlist() models.WatchlistResponse {
	if s.store.HasWatchlist() {
		return s.store.Watchlist()
	}
	return models.WatchlistResponse{Items: []models.WatchlistItem{}}
}

func (s *RealtimeService) Filings() models.FilingsResponse {
	if s.store.HasFilings() {
		return s.store.Filings()
	}
	return models.FilingsResponse{Items: []models.Filing{}}
}

func (s *RealtimeService) Sectors() []models.Sector {
	if s.store.HasSectors() {
		return s.store.Sectors()
	}
	return s.Seeds()
}

func (s *RealtimeService) Health() error {
	if s.finnhub == nil {
		return fmt.Errorf("realtime service not initialized")
	}
	return nil
}