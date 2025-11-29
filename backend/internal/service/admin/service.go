package admin

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/example/ai-avatar-studio/internal/pkg/password"
	"github.com/example/ai-avatar-studio/internal/repository"
)

// Service gives operators CRUD controls on dictionary/model/post data.
type Service struct {
	configs       *repository.ConfigRepository
	roles         *repository.RoleRepository
	community     *repository.CommunityRepository
	users         *repository.UserRepository
	notifications *repository.NotificationRepository
}

func NewService(configs *repository.ConfigRepository, roles *repository.RoleRepository, community *repository.CommunityRepository, users *repository.UserRepository, notifications *repository.NotificationRepository) *Service {
	return &Service{configs: configs, roles: roles, community: community, users: users, notifications: notifications}
}

func (s *Service) Models(ctx context.Context) ([]model.ModelConfig, error) {
	return s.configs.ListModels(ctx, true)
}

func (s *Service) SaveModel(ctx context.Context, payload *model.ModelConfig) (*model.ModelConfig, error) {
	if err := s.configs.SaveModel(ctx, payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func (s *Service) DeleteModel(ctx context.Context, id string) error {
	return s.configs.DeleteModel(ctx, id)
}

func (s *Service) Dictionary(ctx context.Context, group string) ([]model.DictionaryItem, error) {
	return s.configs.ListDictionary(ctx, group)
}

func (s *Service) SaveDictionary(ctx context.Context, item *model.DictionaryItem) (*model.DictionaryItem, error) {
	if err := s.configs.SaveDictionaryItem(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) DeleteDictionary(ctx context.Context, id string) error {
	return s.configs.DeleteDictionaryItem(ctx, id)
}

func (s *Service) Roles(ctx context.Context) ([]model.Role, error) {
	return s.roles.List(ctx, "", 50)
}

func (s *Service) HidePost(ctx context.Context, postID string) error {
	return s.community.UpdateVisibility(ctx, postID, "hidden")
}

func (s *Service) LimitPost(ctx context.Context, postID string) error {
	return s.community.UpdateVisibility(ctx, postID, "limited")
}

func (s *Service) RestorePost(ctx context.Context, postID string) error {
	return s.community.UpdateVisibility(ctx, postID, "public")
}

func (s *Service) ListPosts(ctx context.Context, query, visibility string, limit, offset int) ([]model.CommunityPost, error) {
	return s.community.ListPostsAdmin(ctx, query, visibility, limit, offset)
}

func (s *Service) ListUsers(ctx context.Context, query string, limit, offset int) ([]model.User, error) {
	return s.users.List(ctx, query, limit, offset)
}

func (s *Service) CreateUser(ctx context.Context, username, email, rawPassword, nickname string, isAdmin bool) (*model.User, error) {
	username = strings.TrimSpace(strings.ToLower(username))
	email = strings.TrimSpace(strings.ToLower(email))
	if username == "" || rawPassword == "" {
		return nil, errors.New("username and password required")
	}
	if existing, _ := s.users.FindByUsername(ctx, username); existing != nil {
		return nil, errors.New("username already taken")
	}
	hash, err := password.Hash(rawPassword)
	if err != nil {
		return nil, err
	}
	nickname = strings.TrimSpace(nickname)
	if nickname == "" {
		nickname = username
	}
	if email == "" {
		email = fmt.Sprintf("%s@users.local", username)
	}
	user := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		Nickname:     nickname,
		IsAdmin:      isAdmin,
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) BanUser(ctx context.Context, id string) error {
	return s.users.SetBanStatus(ctx, id, true)
}

func (s *Service) UnbanUser(ctx context.Context, id string) error {
	return s.users.SetBanStatus(ctx, id, false)
}

func (s *Service) DeleteUser(ctx context.Context, id string) error {
	return s.users.SoftDelete(ctx, id)
}

func (s *Service) ListComments(ctx context.Context, postID, visibility string, limit, offset int) ([]model.CommunityComment, error) {
	return s.community.ListCommentsAdmin(ctx, postID, visibility, limit, offset)
}

func (s *Service) HideComment(ctx context.Context, id string) error {
	return s.community.UpdateCommentVisibility(ctx, id, "hidden")
}

func (s *Service) DeleteComment(ctx context.Context, id string) error {
	return s.community.UpdateCommentVisibility(ctx, id, "deleted")
}

// SendNotification creates notifications for specific users.
func (s *Service) SendNotification(ctx context.Context, title, content string, userIDs []string) (int, error) {
	if len(userIDs) == 0 {
		return 0, errors.New("user_ids required")
	}
	created := 0
	for _, uid := range userIDs {
		if strings.TrimSpace(uid) == "" {
			continue
		}
		n := &model.Notification{
			UserID:  uid,
			Type:    "admin",
			Title:   title,
			Content: content,
			IsRead:  false,
		}
		if err := s.notifications.Create(ctx, n); err != nil {
			return created, err
		}
		created++
	}
	return created, nil
}

// BroadcastNotification sends to all users matching query (paginated).
func (s *Service) BroadcastNotification(ctx context.Context, title, content, query string, limit int) (int, error) {
	if limit <= 0 {
		limit = 500
	}
	users, err := s.users.List(ctx, query, limit, 0)
	if err != nil {
		return 0, err
	}
	var ids []string
	for _, u := range users {
		ids = append(ids, u.ID)
	}
	return s.SendNotification(ctx, title, content, ids)
}
