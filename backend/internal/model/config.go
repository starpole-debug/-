package model

import "time"

// ModelConfig describes an LLM option managed by admins.
type ModelConfig struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Provider         string    `json:"provider"`
	BaseURL          string    `json:"base_url"`
	ModelName        string    `json:"model_name"`
	APIKey           string    `json:"-"`
	Temperature      float64   `json:"temperature"`
	MaxTokens        int       `json:"max_tokens"`
	IsDefault        bool      `json:"is_default"`
	IsEnabled        bool      `json:"is_enabled"`
	Status           string    `json:"status"`
	HasAPIKey        bool      `json:"has_api_key"`
	MaxContextTokens int       `json:"max_context_tokens"`
	PriceCoins       int64     `json:"price_coins"`
	PriceHint        string    `json:"price_hint"`
	ShareRolePct     float64   `json:"share_role_pct"`
	SharePresetPct   float64   `json:"share_preset_pct"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// DictionaryItem stores curated vocabulary used by the UI.
type DictionaryItem struct {
	ID          string    `json:"id"`
	Group       string    `json:"group"`
	Key         string    `json:"key"`
	Label       string    `json:"label"`
	Description string    `json:"description"`
	Order       int       `json:"order"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
