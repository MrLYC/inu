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
	"testing"
)

func TestWriteOutput_PrintOnly(t *testing.T) {
	content := "test output content"

	// Test print only (no file output)
	err := WriteOutput(content, true, "")
	if err != nil {
		t.Errorf("WriteOutput failed: %v", err)
	}
}

func TestWriteOutput_FileOnly(t *testing.T) {
	content := "test output content"
	tmpFile, err := os.CreateTemp("", "test-output-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// Test file output only
	err = WriteOutput(content, false, tmpFile.Name())
	if err != nil {
		t.Errorf("WriteOutput failed: %v", err)
	}

	// Verify file content
	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	if string(data) != content {
		t.Errorf("Expected %q, got %q", content, string(data))
	}
}

func TestWriteOutput_Both(t *testing.T) {
	content := "test output content"
	tmpFile, err := os.CreateTemp("", "test-output-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// Test both print and file output
	err = WriteOutput(content, true, tmpFile.Name())
	if err != nil {
		t.Errorf("WriteOutput failed: %v", err)
	}

	// Verify file content
	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	if string(data) != content {
		t.Errorf("Expected %q, got %q", content, string(data))
	}
}

func TestWriteOutput_Neither(t *testing.T) {
	content := "test output content"

	// Test with neither print nor file - should do nothing but not error
	err := WriteOutput(content, false, "")
	if err != nil {
		t.Errorf("WriteOutput failed: %v", err)
	}
}

func TestWriteOutput_InvalidPath(t *testing.T) {
	content := "test output content"

	// Test with invalid file path
	err := WriteOutput(content, false, "/invalid/path/that/does/not/exist/file.txt")
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}
