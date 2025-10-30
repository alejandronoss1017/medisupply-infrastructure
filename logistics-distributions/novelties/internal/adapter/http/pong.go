package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PongHandler responds with a simple JSON pong message. Useful for health
// checks and quick connectivity tests.
func PongHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
