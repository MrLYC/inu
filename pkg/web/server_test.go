package web

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mrlyc/inu/pkg/anonymizer"
)

// mockAnonymizer is a simple mock for testing web routes
type mockAnonymizer struct{}

func (m *mockAnonymizer) AnonymizeText(ctx context.Context, types []string, text string) (string, []*anonymizer.Entity, error) {
	return text, nil, nil
}

func (m *mockAnonymizer) AnonymizeTextStream(ctx context.Context, types []string, text string, writer io.Writer) ([]*anonymizer.Entity, error) {
	return nil, nil
}

func (m *mockAnonymizer) RestoreText(ctx context.Context, entities []*anonymizer.Entity, text string) (string, error) {
	return text, nil
}

// TestStaticFilesEmbedded verifies that static files are properly embedded
func TestStaticFilesEmbedded(t *testing.T) {
	// Verify staticFS is accessible
	entries, err := staticFS.ReadDir("static")
	assert.NoError(t, err, "Should be able to read static directory")
	assert.NotEmpty(t, entries, "Static directory should not be empty")

	// Verify expected files exist
	expectedFiles := []string{
		"static/index.html",
		"static/app.js",
		"static/styles.css",
	}

	for _, filename := range expectedFiles {
		t.Run("File_"+filename, func(t *testing.T) {
			data, err := staticFS.ReadFile(filename)
			assert.NoError(t, err, "File should be embedded: %s", filename)
			assert.NotEmpty(t, data, "File should not be empty: %s", filename)
		})
	}
}

// TestStaticFilesContent verifies basic content of embedded files
func TestStaticFilesContent(t *testing.T) {
	// Test index.html contains expected DOCTYPE
	indexData, err := staticFS.ReadFile("static/index.html")
	assert.NoError(t, err)
	assert.Contains(t, string(indexData), "<!DOCTYPE html>", "index.html should be valid HTML")
	assert.Contains(t, string(indexData), "Inu", "index.html should contain app name")

	// Test styles.css contains CSS rules
	cssData, err := staticFS.ReadFile("static/styles.css")
	assert.NoError(t, err)
	assert.Contains(t, string(cssData), "body", "styles.css should contain CSS rules")

	// Test app.js contains JavaScript
	jsData, err := staticFS.ReadFile("static/app.js")
	assert.NoError(t, err)
	assert.Contains(t, string(jsData), "function", "app.js should contain JavaScript functions")
}

// TestStaticRoutes tests the static file routing behavior
func TestStaticRoutes(t *testing.T) {
	// Create a mock anonymizer
	mockAnon := &mockAnonymizer{}

	// Create server config without authentication for easier testing
	config := &Config{
		Addr:       "127.0.0.1:8080",
		AdminUser:  "",
		AdminToken: "",
	}

	server, err := NewServer(mockAnon, config)
	require.NoError(t, err, "Should create server successfully")

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		contentCheck   func(t *testing.T, body string)
	}{
		{
			name:           "GET /static/app.js returns 200",
			path:           "/static/app.js",
			expectedStatus: http.StatusOK,
			contentCheck: func(t *testing.T, body string) {
				assert.Contains(t, body, "function", "Should return JavaScript content")
			},
		},
		{
			name:           "GET /static/styles.css returns 200",
			path:           "/static/styles.css",
			expectedStatus: http.StatusOK,
			contentCheck: func(t *testing.T, body string) {
				assert.Contains(t, body, "body", "Should return CSS content")
			},
		},
		{
			name:           "GET /static/index.html returns 200",
			path:           "/static/index.html",
			expectedStatus: http.StatusOK,
			contentCheck: func(t *testing.T, body string) {
				assert.Contains(t, body, "<!DOCTYPE html>", "Should return HTML content")
			},
		},
		{
			name:           "GET /static/nonexistent.js returns 404",
			path:           "/static/nonexistent.js",
			expectedStatus: http.StatusNotFound,
			contentCheck:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			server.engine.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code, "Should return expected status code")

			if tt.contentCheck != nil {
				tt.contentCheck(t, w.Body.String())
			}
		})
	}
}

// TestHomePageRoute tests that the home page route works correctly
func TestHomePageRoute(t *testing.T) {
	mockAnon := &mockAnonymizer{}

	config := &Config{
		Addr:       "127.0.0.1:8080",
		AdminUser:  "",
		AdminToken: "",
	}

	server, err := NewServer(mockAnon, config)
	require.NoError(t, err, "Should create server successfully")

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	server.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Should return 200 OK")
	assert.Contains(t, w.Body.String(), "<!DOCTYPE html>", "Should return index.html")
	assert.Contains(t, w.Body.String(), "Inu", "Should contain app name")
}
