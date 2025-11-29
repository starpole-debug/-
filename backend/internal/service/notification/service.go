package notification

import (
	"context"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/example/ai-avatar-studio/internal/repository"
)

// Service surfaces notification list + mark-read flows.
type Service struct {
	repo *repository.NotificationRepository
}

func NewService(repo *repository.NotificationRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, userID string) ([]model.Notification, error) {
	return s.repo.ListByUser(ctx, userID, 50)
}

func (s *Service) MarkRead(ctx context.Context, userID, notificationID string) error {
	if notificationID == "all" {
		return s.repo.MarkAllRead(ctx, userID)
	}
	return s.repo.MarkRead(ctx, userID, notificationID)
}
