package model

import "time"

// ImageProvider represents a drawable model endpoint (e.g., NovelAI).
type ImageProvider struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	BaseURL        string    `json:"base_url"`
	APIKey         string    `json:"api_key"`
	MaxConcurrency int       `json:"max_concurrency"`
	Weight         int       `json:"weight"`
	Status         string    `json:"status"`
	ParamsJSON     string    `json:"params_json,omitempty"`
	SelectedModel  string    `json:"selected_model,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ImagePreset stores JSON-based prompt instructions.
type ImagePreset struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	PresetJSON string    `json:"preset_json"`
	Status     string    `json:"status"`
	PromptModelKey string `json:"prompt_model_key"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ImageJob tracks a single generation request.
type ImageJob struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	SessionID      string    `json:"session_id"`
	MessageID      string    `json:"message_id"`
	ProviderID     string    `json:"provider_id"`
	PresetID       string    `json:"preset_id"`
	Prompt         string    `json:"prompt"`
	NegativePrompt string    `json:"negative_prompt"`
	FinalPrompt    string    `json:"final_prompt"`
	Status         string    `json:"status"`
	ResultURL      string    `json:"result_url"`
	Error          string    `json:"error"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
