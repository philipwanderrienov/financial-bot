package models

type HealthResponse struct {
	Status string `json:"status"`
}

type SummaryResponse struct {
	UpdatedAt string      `json:"updatedAt"`
	Market    MarketQuote `json:"market"`
	Signals   []Signal    `json:"signals"`
}

type MarketQuote struct {
	Symbol        string  `json:"symbol"`
	Price         float64 `json:"price"`
	ChangePercent float64 `json:"changePercent"`
}

type Signal struct {
	Symbol     string  `json:"symbol"`
	Signal     string  `json:"signal"`
	Confidence float64 `json:"confidence"`
	Reason     string  `json:"reason"`
}

type WatchlistResponse struct {
	Items []WatchlistItem `json:"items"`
}

type WatchlistItem struct {
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	ChangePercent float64 `json:"changePercent"`
	Signal        string  `json:"signal"`
}

type FilingsResponse struct {
	Items []Filing `json:"items"`
}

type Filing struct {
	Symbol      string `json:"symbol"`
	Title       string `json:"title"`
	Source      string `json:"source"`
	PublishedAt string `json:"publishedAt"`
	URL         string `json:"url"`
}

type RecommendationResponse struct {
	UpdatedAt  string                 `json:"updatedAt"`
	Symbol     string                 `json:"symbol"`
	Action     string                 `json:"action"`
	Confidence float64                `json:"confidence"`
	Scores     RecommendationScores   `json:"scores"`
	Reasons    []string               `json:"reasons"`
	Sources    RecommendationSources  `json:"sources"`
}

type RecommendationScores struct {
	Technical   float64 `json:"technical"`
	Fundamental float64 `json:"fundamental"`
	News        float64 `json:"news"`
	Risk        float64 `json:"risk"`
}

type RecommendationSources struct {
	MarketData string `json:"marketData"`
	News       string `json:"news"`
	Filings    string `json:"filings"`
}