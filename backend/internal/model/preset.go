package model

import (
	"time"
)

type Preset struct {
	ID          string                 `json:"id" db:"id"`
	CreatorID   string                 `json:"creator_id" db:"creator_id"`
	Name        string                 `json:"name" db:"name"`
	Description string                 `json:"description" db:"description"`
	ModelKey    string                 `json:"model_key" db:"model_key"`
	Blocks      []PresetBlock          `json:"blocks" db:"blocks"`
	GenParams   map[string]interface{} `json:"gen_params" db:"gen_params"`
	IsPublic    bool                   `json:"is_public" db:"is_public"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

type PresetBlock struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Role    string `json:"role"`
	Content string `json:"content"`
	Enabled bool   `json:"enabled"`
	Marker  bool   `json:"marker"`
}
