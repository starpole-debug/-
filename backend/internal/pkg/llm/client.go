package llm

import (
	"context"

	"github.com/example/ai-avatar-studio/internal/model"
)

// Client defines the abstraction for LLM providers.
type Client interface {
	Generate(ctx context.Context, prompt string, model *model.ModelConfig, history []model.ChatMessage) (string, error)
	StreamGenerate(ctx context.Context, prompt string, model *model.ModelConfig, history []model.ChatMessage, onChunk func(contentDelta string, reasoningDelta string)) error
}
