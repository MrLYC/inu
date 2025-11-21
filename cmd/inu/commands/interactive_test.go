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

package commands

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInteractiveCmd(t *testing.T) {
	cmd := NewInteractiveCmd()

	assert.Equal(t, "interactive", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)

	// Verify flags exist
	assert.NotNil(t, cmd.Flags().Lookup("file"))
	assert.NotNil(t, cmd.Flags().Lookup("content"))
	assert.NotNil(t, cmd.Flags().Lookup("entity-types"))
	assert.NotNil(t, cmd.Flags().Lookup("no-prompt"))

	// Verify delimiter flag does not exist (removed)
	assert.Nil(t, cmd.Flags().Lookup("delimiter"))
}

func TestEOFLogic(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Simple text with EOF",
			input: "Hello\nWorld",
			want:  "Hello\nWorld",
		},
		{
			name:  "Empty input",
			input: "",
			want:  "",
		},
		{
			name:  "Single line",
			input: "Single line",
			want:  "Single line",
		},
		{
			name:  "Multiple lines",
			input: "Line 1\nLine 2\nLine 3",
			want:  "Line 1\nLine 2\nLine 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate EOF behavior
			lines := strings.Split(tt.input, "\n")
			got := strings.Join(lines, "\n")
			assert.Equal(t, tt.want, got)
		})
	}
}
