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

	"github.com/mrlyc/inu/pkg/anonymizer"
)

func TestWriteOutput_DefaultPrint(t *testing.T) {
	content := "test output content"

	// Test default behavior: print to stdout (noPrint=false)
	err := WriteOutput(content, false, "")
	if err != nil {
		t.Errorf("WriteOutput failed: %v", err)
	}
}

func TestWriteOutput_NoPrint(t *testing.T) {
	content := "test output content"

	// Test noPrint=true (suppress stdout)
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
	_ = tmpFile.Close()
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	// Test file output with noPrint=true (file only, no stdout)
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

func TestWriteOutput_PrintAndFile(t *testing.T) {
	content := "test output content"
	tmpFile, err := os.CreateTemp("", "test-output-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	_ = tmpFile.Close()
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	// Test both print and file output (noPrint=false)
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

func TestWriteOutput_InvalidPath(t *testing.T) {
	content := "test output content"

	// Test with invalid file path
	err := WriteOutput(content, false, "/invalid/path/that/does/not/exist/file.txt")
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}

func TestWriteEntitiesToStderr_DefaultBehavior(t *testing.T) {
	entities := []*anonymizer.Entity{
		{Key: "email", Values: []string{"user@example.com"}},
		{Key: "phone", Values: []string{"+1234567890"}},
	}

	// Test default behavior: write to stderr (noPrint=false)
	WriteEntitiesToStderr(entities, false)
	// Note: This test only verifies no panic/error occurs
	// Stderr output is not captured in this simple test
}

func TestWriteEntitiesToStderr_NoPrint(t *testing.T) {
	entities := []*anonymizer.Entity{
		{Key: "email", Values: []string{"user@example.com"}},
	}

	// Test noPrint=true (suppress stderr output)
	WriteEntitiesToStderr(entities, true)
	// Should not output anything
}

func TestWriteEntitiesToStderr_EmptyEntities(t *testing.T) {
	// Test with empty entities slice
	WriteEntitiesToStderr([]*anonymizer.Entity{}, false)
	// Should not output anything
}

func TestWriteEntitiesToStderr_NilEntities(t *testing.T) {
	// Test with nil entities
	WriteEntitiesToStderr(nil, false)
	// Should not output anything
}
