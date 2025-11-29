package model

import "time"

// CommunityPost captures long form community content.
type CommunityPost struct {
	ID           string    `json:"id"`
	AuthorID     string    `json:"author_id"`
	AuthorName   string    `json:"author_name"`
	AuthorAvatar string    `json:"author_avatar"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	StyleID      string    `json:"style_id"`
	TopicIDs     []string  `json:"topic_ids"`
	Visibility   string    `json:"visibility"`
	LinkURL      string    `json:"link_url,omitempty"`
	LinkType     string    `json:"link_type,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Attachments  []CommunityAttachment `json:"attachments,omitempty"`
}

// CommunityComment attaches threaded replies to posts.
type CommunityComment struct {
	ID           string    `json:"id"`
	PostID       string    `json:"post_id"`
	AuthorID     string    `json:"author_id"`
	AuthorName   string    `json:"author_name"`
	AuthorAvatar string    `json:"author_avatar"`
	Content      string    `json:"content"`
	Visibility   string    `json:"visibility"`
	CreatedAt    time.Time `json:"created_at"`
}

// CommunityReaction aggregates likes or favorites on posts.
type CommunityReaction struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

// CommunityAttachment stores uploaded assets for posts.
type CommunityAttachment struct {
	ID       string    `json:"id"`
	PostID   string    `json:"post_id"`
	FileURL  string    `json:"file_url"`
	FileType string    `json:"file_type"`
	FileName string    `json:"file_name"`
	FileSize int64     `json:"file_size"`
	CreatedAt time.Time `json:"created_at"`
}
