package handlers

import (
	"context"
	"io"

	"github.com/mrlyc/inu/pkg/anonymizer"
)

// mockAnonymizer is a mock implementation for testing
type mockAnonymizer struct {
	anonymizeFunc func(ctx context.Context, types []string, text string, writer io.Writer) ([]*anonymizer.Entity, error)
	restoreFunc   func(ctx context.Context, entities []*anonymizer.Entity, text string, writer io.Writer) ([]anonymizer.RestoreFailure, error)
}

func (m *mockAnonymizer) Anonymize(ctx context.Context, types []string, text string, writer io.Writer) ([]*anonymizer.Entity, error) {
	if m.anonymizeFunc != nil {
		return m.anonymizeFunc(ctx, types, text, writer)
	}
	return nil, nil
}

func (m *mockAnonymizer) RestoreText(ctx context.Context, entities []*anonymizer.Entity, text string, writer io.Writer) ([]anonymizer.RestoreFailure, error) {
	if m.restoreFunc != nil {
		return m.restoreFunc(ctx, entities, text, writer)
	}
	return nil, nil
}
