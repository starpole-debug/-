package llm

import (
	"context"
	"fmt"

	"github.com/example/ai-avatar-studio/internal/model"
)

// MockClient concatenates prompt + last message to provide a deterministic stub.
type MockClient struct{}

// Generate crafts a naive response that still reflects the user inputs for demos.
func (MockClient) Generate(ctx context.Context, prompt string, model *model.ModelConfig, history []model.ChatMessage) (string, error) {
	_ = ctx
	var latestUser string
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Role == "user" {
			latestUser = history[i].Content
			break
		}
	}
	return fmt.Sprintf("[%s mock] You said: %s\nSystem prompt: %s", model.Name, latestUser, prompt), nil
}

func (MockClient) StreamGenerate(ctx context.Context, prompt string, model *model.ModelConfig, history []model.ChatMessage, onChunk func(contentDelta string, reasoningDelta string)) error {
	resp, _ := MockClient{}.Generate(ctx, prompt, model, history)
	for i := 0; i < len(resp); i += 16 {
		end := i + 16
		if end > len(resp) {
			end = len(resp)
		}
		onChunk(resp[i:end], "")
	}
	return nil
}
