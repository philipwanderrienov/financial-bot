# Finance Agent Contracts

## Project goal
Build a market intelligence app with:
- Go backend
- React + Tailwind frontend
- Data ingestion for market data and issuer filings
- Dashboard for watchlists, market overview, and AI-generated summaries

## API contracts

### `GET /api/health`
Response:
```json
{
  "status": "ok"
}
```

### `GET /api/summary`
Response:
```json
{
  "updatedAt": "2026-04-06T00:00:00Z",
  "market": {
    "symbol": "AAPL",
    "price": 0,
    "changePercent": 0
  },
  "signals": [
    {
      "symbol": "AAPL",
      "signal": "buy",
      "confidence": 0.82,
      "reason": "Earnings momentum remains strong"
    }
  ]
}
```

### `GET /api/watchlist`
Response:
```json
{
  "items": [
    {
      "symbol": "AAPL",
      "name": "Apple Inc.",
      "price": 0,
      "changePercent": 0,
      "signal": "buy"
    }
  ]
}
```

### `GET /api/filings`
Response:
```json
{
  "items": [
    {
      "symbol": "AAPL",
      "title": "Quarterly Report",
      "source": "SEC",
      "publishedAt": "2026-04-06T00:00:00Z",
      "url": "https://example.com"
    }
  ]
}
```

## Frontend UI class contracts
Use these CSS class names for the initial UI:

- `app-shell`
- `app-header`
- `app-title`
- `app-subtitle`
- `dashboard-grid`
- `dashboard-card`
- `dashboard-card__title`
- `dashboard-card__value`
- `dashboard-card__meta`
- `signal-list`
- `signal-item`
- `signal-item__symbol`
- `signal-item__signal`
- `filing-list`
- `filing-item`

## Coding conventions
- Tailwind-first styling in React
- Minimal inline styles
- API responses should be JSON only
- Backend should expose clean REST endpoints
- Keep business logic separate from HTTP handlers
- Use UTC timestamps in API responses
