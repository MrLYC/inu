package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrlyc/inu/pkg/anonymizer"
)

// Anonymizer defines the interface for anonymization operations
type Anonymizer interface {
	AnonymizeText(ctx context.Context, types []string, text string) (string, []*anonymizer.Entity, error)
}

// AnonymizeRequest represents the request body for anonymization
type AnonymizeRequest struct {
	Text        string   `json:"text" binding:"required"`
	EntityTypes []string `json:"entity_types"`
}

// AnonymizeResponse represents the response body for anonymization
type AnonymizeResponse struct {
	AnonymizedText string               `json:"anonymized_text"`
	Entities       []*anonymizer.Entity `json:"entities"`
}

// AnonymizeHandler returns a handler for the anonymize endpoint
func AnonymizeHandler(anon Anonymizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AnonymizeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_input",
				"message": "Invalid JSON format: " + err.Error(),
				"code":    400,
			})
			return
		}

		// Validate text is not empty
		if req.Text == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_input",
				"message": "Text cannot be empty",
				"code":    400,
			})
			return
		}

		// Use default entity types if not specified
		entityTypes := req.EntityTypes
		if len(entityTypes) == 0 {
			entityTypes = anonymizer.DefaultEntityTypes
		}

		// Call anonymizer
		anonymizedText, entities, err := anon.AnonymizeText(c.Request.Context(), entityTypes, req.Text)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "llm_error",
				"message": "Failed to call LLM API: " + err.Error(),
				"code":    500,
			})
			return
		}

		// Return successful response
		c.JSON(http.StatusOK, AnonymizeResponse{
			AnonymizedText: anonymizedText,
			Entities:       entities,
		})
	}
}
