package rag

import (
	"context"
	"strings"

	"github.com/example/ai-avatar-studio/internal/repository"
)

// Service exposes a placeholder retrieval augmented generation helper.
type Service struct {
	documents *repository.DocumentRepository
}

func NewService(documents *repository.DocumentRepository) *Service {
	return &Service{documents: documents}
}

// RetrieveContext returns concatenated snippets to enrich prompts.
func (s *Service) RetrieveContext(ctx context.Context, roleID string) (string, error) {
	chunks, err := s.documents.ListRecentChunks(ctx, 5)
	if err != nil {
		return "", err
	}
	var builder []string
	for _, chunk := range chunks {
		builder = append(builder, chunk.Content)
	}
	return strings.Join(builder, "\n---\n"), nil
}
