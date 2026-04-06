package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	FinnhubAPIKey  string
	RefreshSeconds int
}

func Load() Config {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	refreshSeconds := 15
	if value := os.Getenv("REFRESH_SECONDS"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			refreshSeconds = parsed
		}
	}

	return Config{
		Port:           port,
		FinnhubAPIKey:  os.Getenv("FINNHUB_API_KEY"),
		RefreshSeconds: refreshSeconds,
	}
}