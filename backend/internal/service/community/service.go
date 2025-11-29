package community

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/example/ai-avatar-studio/internal/repository"
	"github.com/example/ai-avatar-studio/internal/task"
)

// Service encapsulates post/comment/reaction flows.
type Service struct {
	repo          *repository.CommunityRepository
	userRepo      *repository.UserRepository
	notifications *repository.NotificationRepository
	configs       *repository.ConfigRepository
	dispatcher    *task.Dispatcher
}

func NewService(repo *repository.CommunityRepository, userRepo *repository.UserRepository, notifications *repository.NotificationRepository, dispatcher *task.Dispatcher, configs *repository.ConfigRepository) *Service {
	return &Service{repo: repo, userRepo: userRepo, notifications: notifications, dispatcher: dispatcher, configs: configs}
}

func (s *Service) Feed(ctx context.Context, sort, filter, search, userID string) ([]model.CommunityPost, error) {
	if filter == "following" && strings.TrimSpace(userID) == "" {
		filter = "all"
	}
	posts, err := s.repo.ListPostsFiltered(ctx, 30, sort, filter, userID, "", search)
	if err != nil && filter == "following" {
		// Fallback: 如果关注列表查询失败（例如表缺失、未初始化），回退到公共流，避免 500。
		return s.repo.ListPostsFiltered(ctx, 30, sort, "all", userID, "", search)
	}
	return posts, err
}

func (s *Service) Detail(ctx context.Context, id, viewerID string) (*model.CommunityPost, []model.CommunityComment, map[string]int, error) {
	post, err := s.repo.FindPost(ctx, id)
	if err != nil || post == nil {
		return nil, nil, nil, errors.New("post not found")
	}
	if post.Visibility != "public" {
		return nil, nil, nil, errors.New("post not found")
	}
	comments, err := s.repo.ListComments(ctx, id)
	if err != nil {
		return nil, nil, nil, err
	}
	counts, err := s.repo.CountReactions(ctx, id)
	if err != nil {
		return nil, nil, nil, err
	}
	if viewerID != "" {
		_ = s.repo.RecordView(ctx, viewerID, post.ID)
	}
	return post, comments, counts, nil
}

func (s *Service) CreatePost(ctx context.Context, authorID string, payload *model.CommunityPost, attachments []string) (*model.CommunityPost, error) {
	if strings.TrimSpace(payload.Title) == "" {
		return nil, errors.New("title required")
	}
	payload.AuthorID = authorID
	if payload.Visibility == "" {
		payload.Visibility = "public"
	}
	linkURL, linkType, err := sanitizeLink(payload.LinkURL, payload.LinkType)
	if err != nil {
		return nil, err
	}
	payload.LinkURL = linkURL
	payload.LinkType = linkType
	if err := s.repo.CreatePost(ctx, payload); err != nil {
		return nil, err
	}
	if len(attachments) > 0 {
		if err := s.repo.CreateAttachments(ctx, payload.ID, attachments); err != nil {
			fmt.Printf("CreateAttachments failed: %v\n", err)
		}
	}
	return payload, nil
}

func (s *Service) Comment(ctx context.Context, userID, postID, content string) (*model.CommunityComment, error) {
	if strings.TrimSpace(content) == "" {
		return nil, errors.New("empty comment")
	}
	comment := &model.CommunityComment{PostID: postID, AuthorID: userID, Content: content}
	if err := s.repo.CreateComment(ctx, comment); err != nil {
		return nil, err
	}
	// Best effort notification to post author.
	post, _ := s.repo.FindPost(ctx, postID)
	if post != nil && post.AuthorID != userID {
		notify := func(ctx context.Context) {
			_ = s.notifications.Create(ctx, &model.Notification{
				UserID:  post.AuthorID,
				Type:    "comment",
				Title:   fmt.Sprintf("New comment on %s", post.Title),
				Content: content,
			})
		}
		if s.dispatcher != nil {
			s.dispatcher.Dispatch(notify)
		} else {
			notify(ctx)
		}
	}
	return comment, nil
}

func (s *Service) React(ctx context.Context, userID, postID, reactionType string) (map[string]int, bool, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, false, errors.New("unauthorized")
	}
	reactionType = strings.ToLower(strings.TrimSpace(reactionType))
	if reactionType == "" {
		reactionType = "like"
	}
	if reactionType != "like" && reactionType != "favorite" {
		return nil, false, errors.New("unsupported reaction")
	}
	activated, err := s.repo.ToggleReaction(ctx, postID, userID, reactionType)
	if err != nil {
		return nil, false, err
	}
	if s.notifications != nil {
		post, _ := s.repo.FindPost(ctx, postID)
		if post != nil && post.AuthorID != userID {
			title := fmt.Sprintf("你的帖子获得新的%s", reactionType)
			if reactionType == "favorite" {
				title = "你的帖子被收藏了"
			} else if reactionType == "like" {
				title = "你的帖子被点赞了"
			}
			content := fmt.Sprintf("%s 刚刚互动了你的帖子《%s》", userID, post.Title)
			notify := func(ctx context.Context) {
				_ = s.notifications.Create(ctx, &model.Notification{
					UserID:  post.AuthorID,
					Type:    reactionType,
					Title:   title,
					Content: content,
				})
			}
			if s.dispatcher != nil {
				s.dispatcher.Dispatch(notify)
			} else {
				notify(ctx)
			}
		}
	}
	counts, _ := s.repo.CountReactions(ctx, postID)
	return counts, activated, nil
}

func (s *Service) Dictionary(ctx context.Context) (map[string][]model.DictionaryItem, error) {
	if s.configs == nil {
		return map[string][]model.DictionaryItem{}, nil
	}
	items, err := s.configs.ListDictionary(ctx, "")
	if err != nil {
		return nil, err
	}
	grouped := make(map[string][]model.DictionaryItem)
	for _, item := range items {
		group := strings.TrimSpace(item.Group)
		if group == "" {
			group = "default"
		}
		grouped[group] = append(grouped[group], item)
	}
	return grouped, nil
}

func (s *Service) Favorites(ctx context.Context, userID string) ([]model.CommunityPost, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("unauthorized")
	}
	return s.repo.ListFavorites(ctx, userID, 30)
}

func (s *Service) RecentViews(ctx context.Context, userID string) ([]model.CommunityPost, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("unauthorized")
	}
	return s.repo.ListRecentViews(ctx, userID, 12)
}

func sanitizeLink(raw string, linkType string) (string, string, error) {
	raw = strings.TrimSpace(raw)
	linkType = strings.TrimSpace(strings.ToLower(linkType))
	if raw == "" {
		return "", linkType, nil
	}
	// prevent javascript: or data: etc.
	lower := strings.ToLower(raw)
	if strings.HasPrefix(lower, "javascript:") || strings.HasPrefix(lower, "data:") {
		return "", "", errors.New("invalid link scheme")
	}
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		if _, err := url.ParseRequestURI(raw); err != nil {
			return "", "", errors.New("invalid link url")
		}
		if linkType == "" {
			linkType = "external"
		}
		return raw, linkType, nil
	}
	// treat as internal path (role/preset)
	if !strings.HasPrefix(raw, "/") {
		raw = "/" + raw
	}
	if linkType == "" {
		linkType = "internal"
	}
	return raw, linkType, nil
}
