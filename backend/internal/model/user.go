package model

import "time"

// User represents an end-user or admin account recorded in the users table.
type User struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Nickname     string     `json:"nickname"`
	AvatarURL    string     `json:"avatar_url"`
	IsAdmin      bool       `json:"is_admin"`
	IsBanned     bool       `json:"is_banned"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	// Optional community stats
	FollowerCount  int  `json:"follower_count,omitempty"`
	FollowingCount int  `json:"following_count,omitempty"`
	IsFollowing    bool `json:"is_following,omitempty"`
}
