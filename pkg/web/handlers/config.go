package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ConfigResponse represents the configuration response
type ConfigResponse struct {
	EntityTypes []string `json:"entity_types"`
}

// ConfigHandler returns a handler for the GET /api/v1/config endpoint
func ConfigHandler(entityTypes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, ConfigResponse{
			EntityTypes: entityTypes,
		})
	}
}
