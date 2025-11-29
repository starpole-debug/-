package memory

import (
	"context"
	"strings"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/example/ai-avatar-studio/internal/repository"
)

// Service represents the long-term preference store per user/role pair.
type Service struct {
	repo *repository.MemoryRepository
}

func NewService(repo *repository.MemoryRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, userID, roleID string) ([]model.MemoryCapsule, error) {
	return s.repo.List(ctx, userID, roleID)
}

func (s *Service) Remember(ctx context.Context, userID, roleID, content string) error {
	if strings.TrimSpace(content) == "" {
		return nil
	}
	capsule := &model.MemoryCapsule{UserID: userID, RoleID: roleID, Content: content}
	return s.repo.Create(ctx, capsule)
}

func (s *Service) Forget(ctx context.Context, userID, capsuleID string) error {
	return s.repo.Delete(ctx, capsuleID, userID)
}
