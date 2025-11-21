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

func TestAnonymizeHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAnon := &mockAnonymizer{
		anonymizeFunc: func(ctx context.Context, types []string, text string) (string, []*anonymizer.Entity, error) {
			return "<个人信息[0].姓名.全名>", []*anonymizer.Entity{
				{
					Key:        "<个人信息[0].姓名.全名>",
					EntityType: "个人信息",
					ID:         "0",
					Category:   "姓名",
					Detail:     "张三",
					Values:     []string{"张三"},
				},
			}, nil
		},
	}

	router := gin.New()
	router.POST("/anonymize", AnonymizeHandler(mockAnon))

	reqBody := AnonymizeRequest{
		Text: "张三的信息",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/anonymize", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response AnonymizeResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response.AnonymizedText != "<个人信息[0].姓名.全名>" {
		t.Errorf("unexpected anonymized text: %s", response.AnonymizedText)
	}

	if len(response.Entities) != 1 {
		t.Errorf("expected 1 entity, got %d", len(response.Entities))
	}
}

func TestAnonymizeHandler_EmptyText(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAnon := &mockAnonymizer{}
	router := gin.New()
	router.POST("/anonymize", AnonymizeHandler(mockAnon))

	reqBody := AnonymizeRequest{
		Text: "",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/anonymize", bytes.NewBuffer(body))
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

func TestAnonymizeHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAnon := &mockAnonymizer{}
	router := gin.New()
	router.POST("/anonymize", AnonymizeHandler(mockAnon))

	req := httptest.NewRequest("POST", "/anonymize", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestAnonymizeHandler_LLMError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAnon := &mockAnonymizer{
		anonymizeFunc: func(ctx context.Context, types []string, text string) (string, []*anonymizer.Entity, error) {
			return "", nil, fmt.Errorf("LLM API error")
		},
	}

	router := gin.New()
	router.POST("/anonymize", AnonymizeHandler(mockAnon))

	reqBody := AnonymizeRequest{
		Text: "some text",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/anonymize", bytes.NewBuffer(body))
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

	if response["error"] != "llm_error" {
		t.Errorf("unexpected error: %v", response["error"])
	}
}

func TestAnonymizeHandler_DefaultEntityTypes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var receivedTypes []string
	mockAnon := &mockAnonymizer{
		anonymizeFunc: func(ctx context.Context, types []string, text string) (string, []*anonymizer.Entity, error) {
			receivedTypes = types
			return "anonymized", []*anonymizer.Entity{}, nil
		},
	}

	router := gin.New()
	router.POST("/anonymize", AnonymizeHandler(mockAnon))

	reqBody := AnonymizeRequest{
		Text: "some text",
		// EntityTypes not specified
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/anonymize", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Should use default entity types
	if len(receivedTypes) != len(anonymizer.DefaultEntityTypes) {
		t.Errorf("expected %d entity types, got %d", len(anonymizer.DefaultEntityTypes), len(receivedTypes))
	}
}

func TestAnonymizeHandler_CustomEntityTypes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var receivedTypes []string
	mockAnon := &mockAnonymizer{
		anonymizeFunc: func(ctx context.Context, types []string, text string) (string, []*anonymizer.Entity, error) {
			receivedTypes = types
			return "anonymized", []*anonymizer.Entity{}, nil
		},
	}

	router := gin.New()
	router.POST("/anonymize", AnonymizeHandler(mockAnon))

	customTypes := []string{"个人信息"}
	reqBody := AnonymizeRequest{
		Text:        "some text",
		EntityTypes: customTypes,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/anonymize", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Should use custom entity types
	if len(receivedTypes) != len(customTypes) {
		t.Errorf("expected %d entity types, got %d", len(customTypes), len(receivedTypes))
	}

	if len(receivedTypes) > 0 && receivedTypes[0] != customTypes[0] {
		t.Errorf("expected entity type '%s', got '%s'", customTypes[0], receivedTypes[0])
	}
}
