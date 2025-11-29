package community

import (
	"context"
	"errors"
	"strings"

	"github.com/example/ai-avatar-studio/internal/model"
)

func (s *Service) GetUserProfile(ctx context.Context, userID, viewerID string) (*model.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	followers, _ := s.repo.CountFollowers(ctx, userID)
	following, _ := s.repo.CountFollowing(ctx, userID)
	isFollowing := false
	if viewer := strings.TrimSpace(viewerID); viewer != "" && viewer != userID {
		isFollowing, _ = s.repo.IsFollowing(ctx, viewer, userID)
	}
	// Sanitize: only return public info
	return &model.User{
		ID:             user.ID,
		Username:       user.Username,
		Nickname:       user.Nickname,
		AvatarURL:      user.AvatarURL,
		IsAdmin:        user.IsAdmin,
		CreatedAt:      user.CreatedAt,
		FollowerCount:  followers,
		FollowingCount: following,
		IsFollowing:    isFollowing,
	}, nil
}

func (s *Service) GetUserPosts(ctx context.Context, userID string) ([]model.CommunityPost, error) {
	return s.repo.ListPostsFiltered(ctx, 50, "latest", "all", "", userID, "")
}

func (s *Service) ToggleFollow(ctx context.Context, followerID, targetID string) (bool, int, int, error) {
	if strings.TrimSpace(followerID) == "" {
		return false, 0, 0, errors.New("unauthorized")
	}
	if strings.TrimSpace(targetID) == "" || followerID == targetID {
		return false, 0, 0, errors.New("invalid target")
	}
	isFollowing, err := s.repo.IsFollowing(ctx, followerID, targetID)
	if err != nil {
		return false, 0, 0, err
	}
	if isFollowing {
		if err := s.repo.UnfollowUser(ctx, followerID, targetID); err != nil {
			return false, 0, 0, err
		}
	} else {
		if err := s.repo.FollowUser(ctx, followerID, targetID); err != nil {
			return false, 0, 0, err
		}
		if s.notifications != nil {
			notify := func(ctx context.Context) {
				_ = s.notifications.Create(ctx, &model.Notification{
					UserID:  targetID,
					Type:    "follow",
					Title:   "有人关注了你",
					Content: "你有新的关注者",
				})
			}
			if s.dispatcher != nil {
				s.dispatcher.Dispatch(notify)
			} else {
				notify(ctx)
			}
		}
	}
	followers, _ := s.repo.CountFollowers(ctx, targetID)
	following, _ := s.repo.CountFollowing(ctx, targetID)
	return !isFollowing, followers, following, nil
}
