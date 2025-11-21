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

func TestSaveAndLoadEntities(t *testing.T) {
	// Create test entities
	testEntities := []*anonymizer.Entity{
		{
			Key:        "<个人信息[0].姓名.全名>",
			EntityType: "个人信息",
			ID:         "0",
			Category:   "姓名",
			Detail:     "张三",
			Values:     []string{"张三"},
		},
		{
			Key:        "<账户信息[0].银行账户.6222021001123456789>",
			EntityType: "账户信息",
			ID:         "0",
			Category:   "银行账户",
			Detail:     "6222021001123456789",
			Values:     []string{"6222021001123456789"},
		},
	}

	// Create temp file
	tmpFile, err := os.CreateTemp("", "test-entities-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	_ = tmpFile.Close()
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	// Test saving entities
	err = SaveEntitiesToYAML(testEntities, tmpFile.Name())
	if err != nil {
		t.Fatalf("SaveEntitiesToYAML failed: %v", err)
	}

	// Test loading entities
	loadedEntities, err := LoadEntitiesFromYAML(tmpFile.Name())
	if err != nil {
		t.Fatalf("LoadEntitiesFromYAML failed: %v", err)
	}

	// Verify loaded entities match original
	if len(loadedEntities) != len(testEntities) {
		t.Fatalf("Expected %d entities, got %d", len(testEntities), len(loadedEntities))
	}

	for i, entity := range loadedEntities {
		expected := testEntities[i]
		if entity.Key != expected.Key {
			t.Errorf("Entity %d: Expected Key %q, got %q", i, expected.Key, entity.Key)
		}
		if entity.EntityType != expected.EntityType {
			t.Errorf("Entity %d: Expected EntityType %q, got %q", i, expected.EntityType, entity.EntityType)
		}
		if entity.ID != expected.ID {
			t.Errorf("Entity %d: Expected ID %s, got %s", i, expected.ID, entity.ID)
		}
		if entity.Category != expected.Category {
			t.Errorf("Entity %d: Expected Category %q, got %q", i, expected.Category, entity.Category)
		}
		if entity.Detail != expected.Detail {
			t.Errorf("Entity %d: Expected Detail %q, got %q", i, expected.Detail, entity.Detail)
		}
	}
}

func TestLoadEntitiesFromYAML_FileNotFound(t *testing.T) {
	_, err := LoadEntitiesFromYAML("/nonexistent/file.yaml")
	if err == nil {
		t.Error("Expected error when file not found, got nil")
	}
}

func TestSaveEntitiesToYAML_InvalidPath(t *testing.T) {
	testEntities := []*anonymizer.Entity{
		{
			Key:        "<test[0].test.test>",
			EntityType: "test",
			ID:         "0",
			Category:   "test",
			Detail:     "test",
			Values:     []string{"test"},
		},
	}

	err := SaveEntitiesToYAML(testEntities, "/invalid/path/that/does/not/exist/file.yaml")
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}

func TestSaveEntitiesToYAML_EmptyEntities(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-entities-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	_ = tmpFile.Close()
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	// Test saving empty entities list
	err = SaveEntitiesToYAML([]*anonymizer.Entity{}, tmpFile.Name())
	if err != nil {
		t.Fatalf("SaveEntitiesToYAML failed: %v", err)
	}

	// Load and verify
	loadedEntities, err := LoadEntitiesFromYAML(tmpFile.Name())
	if err != nil {
		t.Fatalf("LoadEntitiesFromYAML failed: %v", err)
	}

	if len(loadedEntities) != 0 {
		t.Errorf("Expected 0 entities, got %d", len(loadedEntities))
	}
}
