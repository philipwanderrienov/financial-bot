# Backend

Go backend for the Finance Agent app.

## Requirements

- Go 1.22+

## Run

```bash
go mod tidy
go run ./cmd/server
```

The server listens on `:8080` by default. Set `PORT` to override it.

## API endpoints

- `GET /api/health`
- `GET /api/summary`
- `GET /api/watchlist`
- `GET /api/filings`

## Notes

- Responses are JSON only.
- Timestamps are UTC and formatted as RFC3339.
- CORS is enabled for the frontend.