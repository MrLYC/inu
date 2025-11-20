package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// HealthHandler returns a handler for the health check endpoint
func HealthHandler(version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, HealthResponse{
			Status:  "ok",
			Version: version,
		})
	}
}
