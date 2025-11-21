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

package anonymizer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// mockChatModel is a mock implementation of model.BaseChatModel for testing.
// It allows configuring responses without making real LLM API calls.
type mockChatModel struct {
	response      *schema.Message // The mocked response to return
	responseError error           // The error to return (if any)
	streamTokens  []string        // Tokens to stream (for Stream() method)
	streamError   error           // Error to return during streaming
}

// Generate implements model.BaseChatModel.Generate for testing.
func (m *mockChatModel) Generate(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	if m.responseError != nil {
		return nil, m.responseError
	}
	return m.response, nil
}

// BindTools implements model.BaseChatModel.BindTools (not used in tests).
func (m *mockChatModel) BindTools(tools []*schema.ToolInfo) error {
	return nil
}

// Stream implements model.BaseChatModel.Stream for testing.
// Note: For now, this returns a minimal implementation. The streaming functionality
// is tested indirectly through AnonymizeText which uses AnonymizeTextStream internally.
func (m *mockChatModel) Stream(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	if m.streamError != nil {
		return nil, m.streamError
	}
	if len(m.streamTokens) == 0 {
		// Fall back to using Generate response for backward compatibility
		if m.response != nil {
			// Simulate streaming by breaking response into tokens
			tokens := []string{m.response.Content}
			m.streamTokens = tokens
		} else {
			return nil, fmt.Errorf("stream not configured in mock")
		}
	}

	// Create a goroutine-based stream simulator
	// This is a workaround since we can't directly construct schema.StreamReader
	return m.createStreamReaderFromTokens(m.streamTokens)
}

// createStreamReaderFromTokens creates a StreamReader from a token slice
// This is a helper to work around type system limitations in testing
func (m *mockChatModel) createStreamReaderFromTokens(tokens []string) (*schema.StreamReader[*schema.Message], error) {
	// For testing purposes, we'll need to use reflection or accept the limitation
	// that we can't fully test streaming without actual LLM integration.
	// Return an error for now to indicate streaming isn't fully mocked.
	return nil, fmt.Errorf("stream mocking not yet implemented - use Generate() path for testing")
}

// newMockAnonymizeResponse constructs a mock LLM response in the expected format:
// <anonymized_text>
// <<<PAIR>>>
// {"<EntityType[ID].Category.Detail>": ["original_value"]}
func newMockAnonymizeResponse(anonymizedText string, mapping map[string][]string) *schema.Message {
	mappingJSON, _ := json.Marshal(mapping)
	content := fmt.Sprintf("%s\n<<<PAIR>>>\n%s", anonymizedText, string(mappingJSON))
	return &schema.Message{
		Role:    schema.Assistant,
		Content: content,
	}
}

// newMockErrorResponse creates a mock that returns an error.
func newMockErrorResponse(err error) *mockChatModel {
	return &mockChatModel{
		responseError: err,
	}
}

// newMockWithResponse creates a mock that returns a specific response.
func newMockWithResponse(response *schema.Message) *mockChatModel {
	return &mockChatModel{
		response: response,
	}
}
