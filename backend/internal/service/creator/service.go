package creator

import (
	"context"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/example/ai-avatar-studio/internal/repository"
)

// Dashboard aggregates counts displayed in the creator center.
type Dashboard struct {
	TotalRoles   int                  `json:"total_roles"`
	DraftRoles   int                  `json:"draft_roles"`
	Published    int                  `json:"published_roles"`
	RecentEvents []model.RevenueEvent `json:"recent_events"`
	Wallet       *model.CreatorWallet `json:"wallet"`
}

// Service fetches creator specific views.
type Service struct {
	roles   *repository.RoleRepository
	revenue *repository.RevenueRepository
}

func NewService(roles *repository.RoleRepository, revenue *repository.RevenueRepository) *Service {
	return &Service{roles: roles, revenue: revenue}
}

func (s *Service) Dashboard(ctx context.Context, creatorID string) (*Dashboard, error) {
	roles, err := s.roles.ListByCreator(ctx, creatorID)
	if err != nil {
		return nil, err
	}
	wallet, err := s.revenue.GetWallet(ctx, creatorID)
	if err != nil {
		return nil, err
	}
	events, _ := s.revenue.ListEvents(ctx, creatorID, 10)
	dash := &Dashboard{Wallet: wallet, RecentEvents: events}
	for _, role := range roles {
		dash.TotalRoles++
		switch role.Status {
		case "draft":
			dash.DraftRoles++
		case "published":
			dash.Published++
		}
	}
	return dash, nil
}

func (s *Service) Roles(ctx context.Context, creatorID string) ([]model.Role, error) {
	return s.roles.ListByCreator(ctx, creatorID)
}
