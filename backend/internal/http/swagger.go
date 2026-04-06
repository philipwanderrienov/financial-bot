package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const swaggerJSON = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Finance Agent API",
    "version": "1.0.0",
    "description": "API for market intelligence dashboard and mock financial data."
  },
  "servers": [
    {
      "url": "http://localhost:8080"
    }
  ],
  "paths": {
    "/api/health": {
      "get": {
        "summary": "Health check",
        "responses": {
          "200": {
            "description": "Service is healthy"
          }
        }
      }
    },
    "/api/summary": {
      "get": {
        "summary": "Get market summary",
        "responses": {
          "200": {
            "description": "Market summary payload"
          }
        }
      }
    },
    "/api/watchlist": {
      "get": {
        "summary": "Get watchlist",
        "responses": {
          "200": {
            "description": "Watchlist payload"
          }
        }
      }
    },
    "/api/filings": {
      "get": {
        "summary": "Get latest filings",
        "responses": {
          "200": {
            "description": "Filings payload"
          }
        }
      }
    },
    "/api/recommendation": {
      "get": {
        "summary": "Get buy hold sell recommendation",
        "parameters": [
          {
            "name": "symbol",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "example": "AAPL"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Recommendation payload"
          },
          "502": {
            "description": "Recommendation data source error"
          }
        }
      }
    }
  }
}`

func SwaggerJSON(c *gin.Context) {
	c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(swaggerJSON))
}

func SwaggerUI(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Finance Agent API Docs</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
    <style>
      html, body {
        margin: 0;
        padding: 0;
        background: #ffffff;
      }
      #swagger-ui {
        background: #ffffff;
      }
    </style>
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
    <script>
      window.onload = function() {
        window.ui = SwaggerUIBundle({
          url: '/swagger.json',
          dom_id: '#swagger-ui',
          deepLinking: true,
          presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
          layout: 'BaseLayout'
        });
      };
    </script>
  </body>
</html>`)
}
