package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Sectors(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"items": defaultSectors(),
	})
}
