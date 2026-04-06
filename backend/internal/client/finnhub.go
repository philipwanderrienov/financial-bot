package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type FinnhubClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewFinnhubClient(apiKey string) *FinnhubClient {
	return &FinnhubClient{
		baseURL: "https://finnhub.io/api/v1",
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type Quote struct {
	Current  float64 `json:"c"`
	Change   float64 `json:"d"`
	Percent  float64 `json:"dp"`
	High     float64 `json:"h"`
	Low      float64 `json:"l"`
	Open     float64 `json:"o"`
	Previous float64 `json:"pc"`
}

type CompanyProfile struct {
	Name     string `json:"name"`
	Ticker   string `json:"ticker"`
	Currency string `json:"currency"`
	Exchange string `json:"exchange"`
	IPO      string `json:"ipo"`
	WebURL   string `json:"weburl"`
	Country  string `json:"country"`
	Sector   string `json:"finnhubIndustry"`
}

type RecommendationTrend struct {
	Symbol     string `json:"symbol"`
	Buy        int    `json:"buy"`
	Hold       int    `json:"hold"`
	Sell       int    `json:"sell"`
	Period     string `json:"period"`
	StrongBuy  int    `json:"strongBuy"`
	StrongSell int    `json:"strongSell"`
}

type CompanyMetric struct {
	Metric struct {
		MarketCapitalization         float64 `json:"marketCapitalization"`
		PeNormalizedAnnual           float64 `json:"peNormalizedAnnual"`
		ProfitMargin                 float64 `json:"profitMargin"`
		RevenuePerShareTTM           float64 `json:"revenuePerShareTTM"`
		ReturnOnEquityTTM            float64 `json:"returnOnEquityTTM"`
		ReturnOnAssetsTTM            float64 `json:"returnOnAssetsTTM"`
		DividendYieldIndicatedAnnual float64 `json:"dividendYieldIndicatedAnnual"`
		Week52High                   float64 `json:"52WeekHigh"`
		Week52Low                    float64 `json:"52WeekLow"`
		Week52PriceReturnDaily       float64 `json:"52WeekPriceReturnDaily"`
	} `json:"metric"`
}

type NewsArticle struct {
	Headline  string  `json:"headline"`
	Summary   string  `json:"summary"`
	Source    string  `json:"source"`
	URL       string  `json:"url"`
	Datetime  int64   `json:"datetime"`
	Sentiment float64 `json:"sentiment,omitempty"`
}

func (c *FinnhubClient) doJSON(path string, target any) error {
	if c.apiKey == "" {
		return fmt.Errorf("FINNHUB_API_KEY is not set")
	}

	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return err
	}

	q := u.Query()
	q.Set("token", c.apiKey)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("finnhub request failed: %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func (c *FinnhubClient) GetQuote(symbol string) (Quote, error) {
	var out Quote
	err := c.doJSON("/quote?symbol="+url.QueryEscape(symbol), &out)
	return out, err
}

func (c *FinnhubClient) GetProfile(symbol string) (CompanyProfile, error) {
	var out CompanyProfile
	err := c.doJSON("/stock/profile2?symbol="+url.QueryEscape(symbol), &out)
	return out, err
}

func (c *FinnhubClient) GetCompanyNews(symbol string, from string, to string) ([]NewsArticle, error) {
	var out []NewsArticle
	path := fmt.Sprintf("/company-news?symbol=%s&from=%s&to=%s", url.QueryEscape(symbol), url.QueryEscape(from), url.QueryEscape(to))
	err := c.doJSON(path, &out)
	return out, err
}

func (c *FinnhubClient) GetRecommendationTrends(symbol string) ([]RecommendationTrend, error) {
	var out []RecommendationTrend
	err := c.doJSON("/stock/recommendation?symbol="+url.QueryEscape(symbol), &out)
	return out, err
}

func (c *FinnhubClient) GetCompanyMetrics(symbol string) (CompanyMetric, error) {
	var out CompanyMetric
	err := c.doJSON("/stock/metric?symbol="+url.QueryEscape(symbol)+"&metric=all", &out)
	return out, err
}