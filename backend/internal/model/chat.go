package model

import "time"

// ChatSession captures a conversation between a user and a role.
type ChatSession struct {
	ID        string              `json:"id" db:"id"`
	UserID    string              `json:"user_id" db:"user_id"`
	RoleID    string              `json:"role_id" db:"role_id"`
	ModelKey  string              `json:"model_key" db:"model_key"`
	Title     string              `json:"title" db:"title"`
	Summary   string              `json:"summary" db:"summary"` // Auto-generated summary
	LastMsg   string              `json:"last_message,omitempty" db:"last_message"`
	Mode      string              `json:"mode" db:"mode"`       // "sfw", "nsfw"
	Status    string              `json:"status" db:"status"`   // "active", "archived"
	Settings  ChatSessionSettings `json:"settings" db:"settings"`
	CreatedAt time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt time.Time           `json:"updated_at" db:"updated_at"`
}

// ChatMessage stores each user or assistant exchange.
type ChatMessage struct {
	ID          string                 `json:"id"`
	SessionID   string                 `json:"session_id"`
	Role        string                 `json:"role"` // user, assistant, system
	Content     string                 `json:"content"`
	IsImportant bool                   `json:"is_important"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
}

// ChatSessionSettings capture per-session knobs that influence prompting.
type ChatSessionSettings struct {
	Temperature    float64 `json:"temperature"`
	MaxTokens      int     `json:"max_tokens"`
	NarrativeFocus string  `json:"narrative_focus"` // dialogue | balanced | narrative
	ActionRichness string  `json:"action_richness"` // low | medium | high
	SFWMode        bool    `json:"sfw_mode"`
	Immersive      bool    `json:"immersive"`
}

func DefaultChatSessionSettings() ChatSessionSettings {
	return ChatSessionSettings{
		Temperature:    0.7,
		MaxTokens:      512,
		NarrativeFocus: "balanced",
		ActionRichness: "medium",
		SFWMode:        true,
		Immersive:      true,
	}
}
