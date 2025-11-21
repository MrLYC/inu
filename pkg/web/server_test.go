package web

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
