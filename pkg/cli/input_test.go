/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cli

import (
	"os"
	"strings"
	"testing"
)

func TestReadInput_FromFile(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test-input-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	content := "test content from file"
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	_ = tmpFile.Close()

	// Test reading from file (highest priority)
	result, err := ReadInput(tmpFile.Name(), "should be ignored", strings.NewReader("should also be ignored"))
	if err != nil {
		t.Fatalf("ReadInput failed: %v", err)
	}
	if result != content {
		t.Errorf("Expected %q, got %q", content, result)
	}
}

func TestReadInput_FromContent(t *testing.T) {
	content := "test content from string"

	// Test reading from content (second priority)
	result, err := ReadInput("", content, strings.NewReader("should be ignored"))
	if err != nil {
		t.Fatalf("ReadInput failed: %v", err)
	}
	if result != content {
		t.Errorf("Expected %q, got %q", content, result)
	}
}

func TestReadInput_FromStdin(t *testing.T) {
	content := "test content from stdin"
	stdin := strings.NewReader(content)

	// Test reading from stdin (lowest priority)
	result, err := ReadInput("", "", stdin)
	if err != nil {
		t.Fatalf("ReadInput failed: %v", err)
	}
	if result != content {
		t.Errorf("Expected %q, got %q", content, result)
	}
}

func TestReadInput_NoInput(t *testing.T) {
	// Test with no input provided
	_, err := ReadInput("", "", nil)
	if err == nil {
		t.Error("Expected error when no input provided, got nil")
	}
}

func TestReadInput_FileNotFound(t *testing.T) {
	// Test with non-existent file
	_, err := ReadInput("/nonexistent/file.txt", "", nil)
	if err == nil {
		t.Error("Expected error when file not found, got nil")
	}
}

func TestReadInput_EmptyStdin(t *testing.T) {
	// Test with empty stdin
	stdin := strings.NewReader("")
	_, err := ReadInput("", "", stdin)
	if err == nil {
		t.Error("Expected error when stdin is empty, got nil")
	}
}

func TestCheckRequiredEnvVars(t *testing.T) {
	// Save original env vars
	originalAPIKey := os.Getenv("OPENAI_API_KEY")
	originalModelName := os.Getenv("OPENAI_MODEL_NAME")
	defer func() {
		if originalAPIKey != "" {
			_ = os.Setenv("OPENAI_API_KEY", originalAPIKey)
		}
		if originalModelName != "" {
			_ = os.Setenv("OPENAI_MODEL_NAME", originalModelName)
		}
	}()

	tests := []struct {
		name        string
		apiKey      string
		modelName   string
		expectError bool
	}{
		{
			name:        "Both env vars set",
			apiKey:      "test-key",
			modelName:   "gpt-4",
			expectError: false,
		},
		{
			name:        "Missing API key",
			apiKey:      "",
			modelName:   "gpt-4",
			expectError: true,
		},
		{
			name:        "Missing model name",
			apiKey:      "test-key",
			modelName:   "",
			expectError: true,
		},
		{
			name:        "Both missing",
			apiKey:      "",
			modelName:   "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.apiKey != "" {
				_ = os.Setenv("OPENAI_API_KEY", tt.apiKey)
			} else {
				_ = os.Unsetenv("OPENAI_API_KEY")
			}
			if tt.modelName != "" {
				_ = os.Setenv("OPENAI_MODEL_NAME", tt.modelName)
			} else {
				_ = os.Unsetenv("OPENAI_MODEL_NAME")
			}

			err := CheckRequiredEnvVars()
			if tt.expectError && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
