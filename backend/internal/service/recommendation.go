package service

import (
	"fmt"
	"strings"
	"time"

	"finance-agent/backend/internal/client"
	"finance-agent/backend/internal/models"
)

type RecommendationService struct {
	finnhub *client.FinnhubClient
}

func NewRecommendationService(finnhub *client.FinnhubClient) *RecommendationService {
	return &RecommendationService{finnhub: finnhub}
}

func (s *RecommendationService) GetRecommendation(symbol string) (models.RecommendationResponse, error) {
	if symbol == "" {
		symbol = "AAPL"
	}

	quote, profile, news, trends, metrics, err := s.fetchLiveInputs(symbol)
	if err != nil {
		return models.RecommendationResponse{}, err
	}

	technical, technicalReasons := scoreTechnical(quote, trends)
	fundamental, fundamentalReasons := scoreFundamental(profile, metrics, quote)
	newsScore, newsReasons := scoreNews(news)
	risk, riskReasons := scoreRisk(quote, trends, metrics, technical, fundamental, newsScore)

	confidence := clamp((technical*0.32)+(fundamental*0.30)+(newsScore*0.18)+(100-risk)*0.20, 0, 100)
	action := decideAction(confidence, technical, fundamental, newsScore, risk)

	reasons := buildReasons(symbol, technicalReasons, fundamentalReasons, newsReasons, riskReasons, news)

	return models.RecommendationResponse{
		UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
		Symbol:     symbol,
		Action:     action,
		Confidence: confidence,
		Scores: models.RecommendationScores{
			Technical:   technical,
			Fundamental: fundamental,
			News:        newsScore,
			Risk:        risk,
		},
		Reasons: reasons,
		Sources: models.RecommendationSources{
			MarketData: "Finnhub quote + recommendation trend + metric APIs",
			News:       "Finnhub company news API",
			Filings:    "SEC EDGAR / future integration",
		},
	}, nil
}

func (s *RecommendationService) fetchLiveInputs(symbol string) (client.Quote, client.CompanyProfile, []client.NewsArticle, []client.RecommendationTrend, client.CompanyMetric, error) {
	if s.finnhub == nil {
		return client.Quote{}, client.CompanyProfile{}, nil, nil, client.CompanyMetric{}, fmt.Errorf("finnhub client not configured")
	}

	quote, err := s.finnhub.GetQuote(symbol)
	if err != nil {
		return client.Quote{}, client.CompanyProfile{}, nil, nil, client.CompanyMetric{}, err
	}

	profile, err := s.finnhub.GetProfile(symbol)
	if err != nil {
		profile = client.CompanyProfile{}
	}

	end := time.Now().UTC()
	start := end.AddDate(0, 0, -7)

	news, err := s.finnhub.GetCompanyNews(symbol, start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		news = []client.NewsArticle{}
	}

	trends, err := s.finnhub.GetRecommendationTrends(symbol)
	if err != nil {
		trends = []client.RecommendationTrend{}
	}

	metrics, err := s.finnhub.GetCompanyMetrics(symbol)
	if err != nil {
		metrics = client.CompanyMetric{}
	}

	return quote, profile, news, trends, metrics, nil
}

// Score formula notes:
// Technical: starts at 50 and moves with live price momentum, intraday range, and analyst trend balance.
// Fundamental: starts at 50 and adjusts for live profile completeness plus valuation/profitability metrics.
// News: starts at 50 and adjusts by a simple keyword sentiment pass across recent headlines.
// Risk: starts at 50 and increases with volatility, bearish analyst balance, and weak fundamentals.
// All helpers clamp to 0-100 so the frontend always receives stable score bounds.
func scoreTechnical(quote client.Quote, trends []client.RecommendationTrend) (float64, []string) {
	score := 50.0
	reasons := []string{}

	if quote.Current > 0 && quote.Previous > 0 {
		changePct := ((quote.Current - quote.Previous) / quote.Previous) * 100
		score += clamp(changePct*4, -18, 18)
		reasons = append(reasons, fmt.Sprintf("Price change vs previous close is %.2f%%", changePct))
	}

	if quote.High > 0 && quote.Low > 0 && quote.Current > 0 {
		rangePct := ((quote.High - quote.Low) / quote.Current) * 100
		score += clamp(12-(rangePct*1.2), -12, 12)
		reasons = append(reasons, fmt.Sprintf("Intraday range is %.2f%% of price", rangePct))
	}

	if len(trends) > 0 {
		t := trends[0]
		bullish := float64(t.StrongBuy*2 + t.Buy)
		bearish := float64(t.StrongSell*2 + t.Sell)
		balance := bullish - bearish
		score += clamp(balance*3, -15, 15)
		reasons = append(reasons, fmt.Sprintf("Analyst trend balance for %s is %+0.f", t.Period, balance))
	}

	if quote.Current <= 0 {
		score = 35
		reasons = append(reasons, "Quote is unavailable or invalid")
	}

	return clamp(score, 0, 100), reasons
}

func scoreFundamental(profile client.CompanyProfile, metrics client.CompanyMetric, quote client.Quote) (float64, []string) {
	score := 50.0
	reasons := []string{}

	if profile.Name != "" {
		score += 4
		reasons = append(reasons, "Company profile is available")
	}
	if profile.Exchange != "" {
		score += 3
		reasons = append(reasons, "Exchange metadata is available")
	}
	if profile.Sector != "" {
		score += 3
		reasons = append(reasons, "Sector metadata is available")
	}
	if profile.Country == "US" {
		score += 2
	}

	if metrics.Metric.MarketCapitalization > 0 {
		score += 4
		reasons = append(reasons, fmt.Sprintf("Market cap data is available at %.0f", metrics.Metric.MarketCapitalization))
	}
	if metrics.Metric.ProfitMargin > 0 {
		score += clamp(metrics.Metric.ProfitMargin*40, 0, 8)
		reasons = append(reasons, fmt.Sprintf("Profit margin is %.2f", metrics.Metric.ProfitMargin))
	}
	if metrics.Metric.ReturnOnEquityTTM > 0 {
		score += clamp(metrics.Metric.ReturnOnEquityTTM*25, 0, 8)
		reasons = append(reasons, fmt.Sprintf("Return on equity is %.2f", metrics.Metric.ReturnOnEquityTTM))
	}
	if metrics.Metric.PeNormalizedAnnual > 0 && metrics.Metric.PeNormalizedAnnual < 30 {
		score += 4
		reasons = append(reasons, fmt.Sprintf("Normalized P/E is %.2f", metrics.Metric.PeNormalizedAnnual))
	}
	if metrics.Metric.DividendYieldIndicatedAnnual > 0 {
		score += 1
	}
	if quote.Current > 0 && metrics.Metric.Week52High > 0 && metrics.Metric.Week52Low > 0 {
		nearHigh := (quote.Current - metrics.Metric.Week52Low) / (metrics.Metric.Week52High - metrics.Metric.Week52Low)
		score += clamp((nearHigh*10)-5, -5, 5)
		reasons = append(reasons, "Price is positioned relative to the 52-week range")
	}

	return clamp(score, 0, 100), reasons
}

func scoreNews(news []client.NewsArticle) (float64, []string) {
	if len(news) == 0 {
		return 50, []string{"No recent company news was returned by Finnhub"}
	}

	score := 50.0
	reasons := []string{fmt.Sprintf("%d recent headlines were returned", len(news))}

	for _, item := range news {
		text := fmt.Sprintf("%s %s", item.Headline, item.Summary)
		score += sentimentFromText(text)
	}

	score += clamp(10-float64(len(news)), -6, 6)

	return clamp(score, 0, 100), reasons
}

func scoreRisk(quote client.Quote, trends []client.RecommendationTrend, metrics client.CompanyMetric, technical, fundamental, newsScore float64) (float64, []string) {
	risk := 50.0
	reasons := []string{}

	if quote.Current > 0 && quote.High > 0 && quote.Low > 0 {
		intradayVolatility := ((quote.High - quote.Low) / quote.Current) * 100
		risk += clamp(intradayVolatility*2, 0, 18)
		reasons = append(reasons, fmt.Sprintf("Intraday volatility is %.2f%%", intradayVolatility))
	}

	if len(trends) > 0 {
		t := trends[0]
		bearish := float64(t.Sell + (t.StrongSell * 2))
		bullish := float64(t.Buy + (t.StrongBuy * 2))
		if bearish > bullish {
			risk += clamp((bearish-bullish)*3, 0, 15)
			reasons = append(reasons, "Analyst trends lean bearish")
		}
	}

	if metrics.Metric.ProfitMargin > 0 {
		risk += clamp((1-metrics.Metric.ProfitMargin)*12, 0, 10)
	}
	if metrics.Metric.ReturnOnAssetsTTM > 0 {
		risk += clamp((1-metrics.Metric.ReturnOnAssetsTTM)*8, 0, 8)
	}

	risk += clamp(55-technical, -8, 8)
	risk += clamp(55-fundamental, -8, 8)
	risk -= clamp(newsScore-50, -5, 5)

	if quote.Current <= 0 {
		risk = 75
		reasons = append(reasons, "Quote is unavailable or invalid")
	}

	return clamp(risk, 0, 100), reasons
}

func sentimentFromText(text string) float64 {
	lower := normalizeText(text)
	score := 0.0

	positive := []string{"beat", "growth", "strong", "surge", "profit", "upgrade", "bull", "expand", "record", "raise"}
	negative := []string{"miss", "weak", "drop", "loss", "downgrade", "bear", "decline", "lawsuit", "cut", "guidance"}

	for _, token := range positive {
		if contains(lower, token) {
			score += 3
		}
	}
	for _, token := range negative {
		if contains(lower, token) {
			score -= 3
		}
	}
	return score
}

func decideAction(confidence, technical, fundamental, newsScore, risk float64) string {
	if confidence >= 65 && technical >= 50 && fundamental >= 45 && newsScore >= 45 && risk <= 55 {
		return "buy"
	}
	if confidence <= 40 || risk >= 70 {
		return "sell"
	}
	return "hold"
}

func buildReasons(symbol string, technicalReasons, fundamentalReasons, newsReasons, riskReasons []string, news []client.NewsArticle) []string {
	reasons := []string{
		fmt.Sprintf("%s recommendation is based on live Finnhub quote, trend, profile, metric, and news data", symbol),
	}

	reasons = append(reasons, technicalReasons...)
	reasons = append(reasons, fundamentalReasons...)
	reasons = append(reasons, newsReasons...)
	reasons = append(reasons, riskReasons...)

	if len(news) > 0 {
		reasons = append(reasons, fmt.Sprintf("Latest headline: %s", news[0].Headline))
	}

	return reasons
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func normalizeText(s string) string {
	out := strings.Builder{}
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			out.WriteRune(r + 32)
		} else {
			out.WriteRune(r)
		}
	}
	return out.String()
}

func contains(text, sub string) bool {
	return len(sub) > 0 && len(text) >= len(sub) && (indexOf(text, sub) >= 0)
}

func indexOf(text, sub string) int {
	for i := 0; i+len(sub) <= len(text); i++ {
		if text[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
