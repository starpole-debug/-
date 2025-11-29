package store

import (
	"context"
	"errors"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/example/ai-avatar-studio/internal/repository"
	"github.com/example/ai-avatar-studio/internal/service/revenue"
)

// Service handles monetisation entry points from the store/tipping UI.
type Service struct {
	roles   *repository.RoleRepository
	revenue *revenue.Service
}

func NewService(roles *repository.RoleRepository, revenue *revenue.Service) *Service {
	return &Service{roles: roles, revenue: revenue}
}

func (s *Service) TipRole(ctx context.Context, userID, roleID string, amount int64) (*model.RevenueEvent, *model.CreatorWallet, error) {
	role, err := s.roles.FindByID(ctx, roleID)
	if err != nil || role == nil {
		return nil, nil, errors.New("role not found")
	}
	return s.revenue.RecordEvent(ctx, role.CreatorID, userID, roleID, "tip", amount)
}
