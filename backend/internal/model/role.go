package model

import "time"

// Role describes an AI persona created by community creators.
type Role struct {
	ID          string                 `json:"id"`
	CreatorID   string                 `json:"creator_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	AvatarURL   string                 `json:"avatar_url"`
	Tags        []string               `json:"tags"`
	Abilities   []string               `json:"abilities"`
	AllowClone  bool                   `json:"allow_clone"`
	Status      string                 `json:"status"`
	Version     string                 `json:"role_version"`
	Data        map[string]interface{} `json:"data,omitempty"`
	FavoriteCnt int                    `json:"favorite_count,omitempty"`
	IsFavorited bool                   `json:"is_favorited,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// RoleVersion snapshots long form prompt templates for auditability.
type RoleVersion struct {
	ID        string    `json:"id"`
	RoleID    string    `json:"role_id"`
	Prompt    string    `json:"prompt"`
	CreatedAt time.Time `json:"created_at"`
}
