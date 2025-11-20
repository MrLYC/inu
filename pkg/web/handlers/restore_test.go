package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mrlyc/inu/pkg/anonymizer"
)

func TestRestoreHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAnon := &mockAnonymizer{
		restoreFunc: func(ctx context.Context, entities []*anonymizer.Entity, text string) (string, error) {
			return "张三", nil
		},
	}

	router := gin.New()
	router.POST("/restore", RestoreHandler(mockAnon))

	reqBody := RestoreRequest{
		AnonymizedText: "<个人信息[0].姓名.张三>",
		Entities: []*anonymizer.Entity{
			{
				Key:    "<个人信息[0].姓名.张三>",
				Values: []string{"张三"},
			},
		},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/restore", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response RestoreResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response.RestoredText != "张三" {
		t.Errorf("unexpected restored text: %s", response.RestoredText)
	}
}

func TestRestoreHandler_EmptyAnonymizedText(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAnon := &mockAnonymizer{}
	router := gin.New()
	router.POST("/restore", RestoreHandler(mockAnon))

	reqBody := RestoreRequest{
		AnonymizedText: "",
		Entities:       []*anonymizer.Entity{},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/restore", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response["error"] != "invalid_input" {
		t.Errorf("unexpected error: %v", response["error"])
	}
}

func TestRestoreHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAnon := &mockAnonymizer{}
	router := gin.New()
	router.POST("/restore", RestoreHandler(mockAnon))

	req := httptest.NewRequest("POST", "/restore", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestRestoreHandler_RestoreError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAnon := &mockAnonymizer{
		restoreFunc: func(ctx context.Context, entities []*anonymizer.Entity, text string) (string, error) {
			return "", fmt.Errorf("restore error")
		},
	}

	router := gin.New()
	router.POST("/restore", RestoreHandler(mockAnon))

	reqBody := RestoreRequest{
		AnonymizedText: "some text",
		Entities:       []*anonymizer.Entity{},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/restore", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response["error"] != "restore_error" {
		t.Errorf("unexpected error: %v", response["error"])
	}
}

func TestRestoreHandler_EmptyEntities(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAnon := &mockAnonymizer{
		restoreFunc: func(ctx context.Context, entities []*anonymizer.Entity, text string) (string, error) {
			// Empty entities should still work, just return original text
			return text, nil
		},
	}

	router := gin.New()
	router.POST("/restore", RestoreHandler(mockAnon))

	reqBody := RestoreRequest{
		AnonymizedText: "plain text",
		Entities:       []*anonymizer.Entity{},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/restore", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response RestoreResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response.RestoredText != "plain text" {
		t.Errorf("unexpected restored text: %s", response.RestoredText)
	}
}
