package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"finance-agent/backend/internal/config"
	httpserver "finance-agent/backend/internal/http"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	cfg := config.Load()
	router := httpserver.SetupRouter(cfg)

	log.Printf("Backend running at http://localhost:%s", cfg.Port)
	log.Printf("Swagger UI: http://localhost:%s/swagger", cfg.Port)
	log.Printf("Swagger JSON: http://localhost:%s/swagger.json", cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
