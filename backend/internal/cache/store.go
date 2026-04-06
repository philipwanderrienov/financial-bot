package cache

import (
	"sync"
	"time"

	"finance-agent/backend/internal/models"
)

type SnapshotStore struct {
	mu sync.RWMutex

	updatedAt    time.Time
	summary      models.SummaryResponse
	watchlist    models.WatchlistResponse
	filings      models.FilingsResponse
	recommend    map[string]models.RecommendationResponse
	sectors      []models.Sector
	hasSummary   bool
	hasWatchlist bool
	hasFilings   bool
	hasSectors   bool
}

func NewSnapshotStore() *SnapshotStore {
	return &SnapshotStore{
		recommend: make(map[string]models.RecommendationResponse),
	}
}

func (s *SnapshotStore) UpdateSummary(summary models.SummaryResponse) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.summary = summary
	s.hasSummary = true
	s.updatedAt = time.Now().UTC()
}

func (s *SnapshotStore) UpdateWatchlist(watchlist models.WatchlistResponse) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.watchlist = watchlist
	s.hasWatchlist = true
	s.updatedAt = time.Now().UTC()
}

func (s *SnapshotStore) UpdateFilings(filings models.FilingsResponse) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.filings = filings
	s.hasFilings = true
	s.updatedAt = time.Now().UTC()
}

func (s *SnapshotStore) UpdateRecommendation(symbol string, rec models.RecommendationResponse) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.recommend == nil {
		s.recommend = make(map[string]models.RecommendationResponse)
	}
	s.recommend[symbol] = rec
}

func (s *SnapshotStore) UpdateSectors(sectors []models.Sector) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sectors = cloneSectors(sectors)
	s.hasSectors = true
	s.updatedAt = time.Now().UTC()
}

func (s *SnapshotStore) Summary() models.SummaryResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.summary
}

func (s *SnapshotStore) Watchlist() models.WatchlistResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.watchlist
}

func (s *SnapshotStore) Filings() models.FilingsResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.filings
}

func (s *SnapshotStore) Recommendation(symbol string) (models.RecommendationResponse, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rec, ok := s.recommend[symbol]
	return rec, ok
}

func (s *SnapshotStore) Sectors() []models.Sector {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return cloneSectors(s.sectors)
}

func (s *SnapshotStore) HasSummary() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.hasSummary
}

func (s *SnapshotStore) HasWatchlist() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.hasWatchlist
}

func (s *SnapshotStore) HasFilings() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.hasFilings
}

func (s *SnapshotStore) HasSectors() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.hasSectors
}

func cloneSectors(sectors []models.Sector) []models.Sector {
	out := make([]models.Sector, len(sectors))
	for i, sector := range sectors {
		symbols := make([]string, len(sector.Symbols))
		copy(symbols, sector.Symbols)
		out[i] = models.Sector{
			Key:     sector.Key,
			Label:   sector.Label,
			Symbols: symbols,
		}
	}
	return out
}