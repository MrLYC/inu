package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestConfigHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		entityTypes []string
		wantStatus  int
		wantTypes   []string
	}{
		{
			name:        "returns default entity types",
			entityTypes: []string{"PERSON", "ORG", "EMAIL"},
			wantStatus:  http.StatusOK,
			wantTypes:   []string{"PERSON", "ORG", "EMAIL"},
		},
		{
			name:        "returns empty list when no types configured",
			entityTypes: []string{},
			wantStatus:  http.StatusOK,
			wantTypes:   []string{},
		},
		{
			name:        "returns custom entity types",
			entityTypes: []string{"CUSTOM1", "CUSTOM2"},
			wantStatus:  http.StatusOK,
			wantTypes:   []string{"CUSTOM1", "CUSTOM2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/api/v1/config", ConfigHandler(tt.entityTypes))

			req := httptest.NewRequest(http.MethodGet, "/api/v1/config", nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)

			var resp ConfigResponse
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantTypes, resp.EntityTypes)
		})
	}
}
