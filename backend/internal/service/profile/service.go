package profile

import (
	"context"
	"errors"
	"strings"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/example/ai-avatar-studio/internal/repository"
)

// Service aggregates user-centric data for the personal homepage.
type Service struct {
	users     *repository.UserRepository
	community *repository.CommunityRepository
	chats     *repository.ChatRepository
	roles     *repository.RoleRepository
	assets    *repository.UserAssetRepository
	revenue   *repository.RevenueRepository
}

func NewService(users *repository.UserRepository, community *repository.CommunityRepository, chats *repository.ChatRepository, roles *repository.RoleRepository, assets *repository.UserAssetRepository, revenue *repository.RevenueRepository) *Service {
	return &Service{
		users:     users,
		community: community,
		chats:     chats,
		roles:     roles,
		assets:    assets,
		revenue:   revenue,
	}
}

// Overview returns the aggregated user homepage payload.
func (s *Service) Overview(ctx context.Context, userID string) (*model.UserHome, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("unauthorized")
	}
	user, err := s.users.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	assets, err := s.assets.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	favorites, err := s.community.ListFavorites(ctx, userID, 6)
	if err != nil {
		return nil, err
	}
	favCount, err := s.community.CountFavorites(ctx, userID)
	if err != nil {
		return nil, err
	}
	recentViews, err := s.community.ListRecentViews(ctx, userID, 6)
	if err != nil {
		return nil, err
	}
	myPosts, err := s.community.ListPostsByAuthor(ctx, userID, 5)
	if err != nil {
		return nil, err
	}
	postTotal, err := s.community.CountPostsByAuthor(ctx, userID)
	if err != nil {
		return nil, err
	}

	sessionTotal, err := s.chats.CountSessions(ctx, userID)
	if err != nil {
		sessionTotal = 0
	}

	var rolesTotal, publishedRoles int
	if s.roles != nil {
		rolesTotal, publishedRoles, err = s.roles.CountByCreator(ctx, userID)
		if err != nil {
			rolesTotal, publishedRoles = 0, 0
		}
	}

	var walletBalance int64
	if s.revenue != nil {
		if wallet, err := s.revenue.GetWallet(ctx, userID); err == nil && wallet != nil {
			walletBalance = wallet.AvailableBalance
		}
	}

	return &model.UserHome{
		User: &model.User{
			ID:        user.ID,
			Username:  user.Username,
			Nickname:  user.Nickname,
			AvatarURL: user.AvatarURL,
			CreatedAt: user.CreatedAt,
		},
		Assets: assets,
		Stats: model.UserHomeStats{
			SessionTotal:    sessionTotal,
			FavoriteTotal:   favCount,
			RecentViewTotal: len(recentViews),
			PostTotal:       postTotal,
		},
		Favorites:   favorites,
		RecentViews: recentViews,
		MyPosts:     myPosts,
		CreatorHint: model.CreatorHomePrompt{
			RolesTotal:       rolesTotal,
			PublishedRoles:   publishedRoles,
			WalletBalance:    walletBalance,
			HasCreatorAccess: rolesTotal > 0,
		},
	}, nil
}
