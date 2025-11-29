package role

import (
	"context"
	"errors"
	"strings"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/example/ai-avatar-studio/internal/repository"
)

// Service keeps the business rules around roles and publishing workflow.
type Service struct {
	roles *repository.RoleRepository
}

func NewService(roles *repository.RoleRepository) *Service {
	return &Service{roles: roles}
}

func (s *Service) List(ctx context.Context) ([]model.Role, error) {
	return s.roles.List(ctx, "published", 30)
}

func (s *Service) Featured(ctx context.Context) ([]model.Role, error) {
	return s.roles.List(ctx, "published", 12)
}

func (s *Service) Get(ctx context.Context, id string) (*model.Role, error) {
	return s.roles.FindByID(ctx, id)
}

func (s *Service) Save(ctx context.Context, creatorID string, payload *model.Role) (*model.Role, error) {
	if strings.TrimSpace(payload.Name) == "" {
		return nil, errors.New("name required")
	}
	if payload.ID != "" {
		existing, err := s.roles.FindByID(ctx, payload.ID)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.CreatorID != creatorID {
			return nil, errors.New("forbidden")
		}
	}
	payload.CreatorID = creatorID
	if payload.Status == "" {
		payload.Status = "draft"
	}
	if err := s.roles.Save(ctx, payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func (s *Service) Publish(ctx context.Context, roleID string) error {
	return s.roles.UpdateStatus(ctx, roleID, "published")
}

func (s *Service) Archive(ctx context.Context, roleID string) error {
	return s.roles.UpdateStatus(ctx, roleID, "archived")
}

func (s *Service) ListByCreator(ctx context.Context, creatorID string) ([]model.Role, error) {
	return s.roles.ListByCreator(ctx, creatorID)
}

func (s *Service) SnapshotPrompt(ctx context.Context, roleID, prompt string) error {
	return s.roles.CreateVersion(ctx, roleID, prompt)
}

func (s *Service) Favorite(ctx context.Context, userID, roleID string) error {
	if strings.TrimSpace(roleID) == "" {
		return errors.New("role id required")
	}
	role, err := s.roles.FindByID(ctx, roleID)
	if err != nil || role == nil {
		return errors.New("role not found")
	}
	return s.roles.Favorite(ctx, userID, roleID)
}

func (s *Service) Unfavorite(ctx context.Context, userID, roleID string) error {
	if strings.TrimSpace(roleID) == "" {
		return errors.New("role id required")
	}
	return s.roles.Unfavorite(ctx, userID, roleID)
}

func (s *Service) ListFavorites(ctx context.Context, userID string) ([]model.Role, error) {
	return s.roles.ListFavorites(ctx, userID, 50)
}

func (s *Service) PopulateFavoriteMetadata(ctx context.Context, userID string, role *model.Role) {
	if role == nil {
		return
	}
	if cnt, err := s.roles.CountFavorites(ctx, role.ID); err == nil {
		role.FavoriteCnt = cnt
	}
	if userID == "" {
		return
	}
	if liked, err := s.roles.IsFavorited(ctx, userID, role.ID); err == nil {
		role.IsFavorited = liked
	}
}
