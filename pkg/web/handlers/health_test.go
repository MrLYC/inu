package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/health", HealthHandler("v0.1.0"))

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response HealthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response.Status != "ok" {
		t.Errorf("expected status 'ok', got '%s'", response.Status)
	}

	if response.Version != "v0.1.0" {
		t.Errorf("expected version 'v0.1.0', got '%s'", response.Version)
	}
}
