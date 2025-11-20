package handlers

import (
	"context"

	"github.com/mrlyc/inu/pkg/anonymizer"
)

// mockAnonymizer is a mock implementation for testing
type mockAnonymizer struct {
	anonymizeFunc func(ctx context.Context, types []string, text string) (string, []*anonymizer.Entity, error)
	restoreFunc   func(ctx context.Context, entities []*anonymizer.Entity, text string) (string, error)
}

func (m *mockAnonymizer) AnonymizeText(ctx context.Context, types []string, text string) (string, []*anonymizer.Entity, error) {
	if m.anonymizeFunc != nil {
		return m.anonymizeFunc(ctx, types, text)
	}
	return "", nil, nil
}

func (m *mockAnonymizer) RestoreText(ctx context.Context, entities []*anonymizer.Entity, text string) (string, error) {
	if m.restoreFunc != nil {
		return m.restoreFunc(ctx, entities, text)
	}
	return "", nil
}
