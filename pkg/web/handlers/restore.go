package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mrlyc/inu/pkg/anonymizer"
)

// Restorer defines the interface for restoration operations
type Restorer interface {
	RestoreText(ctx context.Context, entities []*anonymizer.Entity, text string) (string, error)
}

// RestoreRequest represents the request body for restoration
type RestoreRequest struct {
	AnonymizedText string               `json:"anonymized_text" binding:"required"`
	Entities       []*anonymizer.Entity `json:"entities" binding:"required"`
}

// RestoreResponse represents the response body for restoration
type RestoreResponse struct {
	RestoredText string `json:"restored_text"`
}

// RestoreHandler returns a handler for the restore endpoint
func RestoreHandler(anon Restorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RestoreRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_input",
				"message": "Invalid JSON format: " + err.Error(),
				"code":    400,
			})
			return
		}

		// Validate anonymized text is not empty
		if req.AnonymizedText == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_input",
				"message": "Anonymized text cannot be empty",
				"code":    400,
			})
			return
		}

		// Call anonymizer to restore text
		restoredText, err := anon.RestoreText(c.Request.Context(), req.Entities, req.AnonymizedText)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "restore_error",
				"message": "Failed to restore text: " + err.Error(),
				"code":    500,
			})
			return
		}

		// Return successful response
		c.JSON(http.StatusOK, RestoreResponse{
			RestoredText: restoredText,
		})
	}
}
