package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/example/ai-avatar-studio/internal/model"
)

const defaultAPIBase = "https://api.openai.com/v1"

// RouterClient dispatches to either the HTTP client or mock fallback based on model provider.
type RouterClient struct {
	http *HTTPClient
	mock Client
}

// NewRouterClient builds a multi-provider client that supports HTTP + mock providers.
func NewRouterClient(httpClient *http.Client) *RouterClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 150 * time.Second}
	}
	return &RouterClient{
		http: NewHTTPClient(httpClient),
		mock: MockClient{},
	}
}

// Generate delegates to the proper provider implementation.
func (r *RouterClient) Generate(ctx context.Context, prompt string, cfg *model.ModelConfig, history []model.ChatMessage) (string, error) {
	if cfg == nil {
		return "", errors.New("model not configured")
	}
	switch strings.ToLower(strings.TrimSpace(cfg.Provider)) {
	case "mock":
		return r.mock.Generate(ctx, prompt, cfg, history)
	default:
		if cfg.APIKey == "" {
			log.Printf("llm: missing api key for model %s, falling back to mock", cfg.ID)
			return r.mock.Generate(ctx, prompt, cfg, history)
		}
		return r.http.Generate(ctx, prompt, cfg, history)
	}
}

func (r *RouterClient) StreamGenerate(ctx context.Context, prompt string, cfg *model.ModelConfig, history []model.ChatMessage, onChunk func(contentDelta string, reasoningDelta string)) error {
	if cfg == nil {
		return errors.New("model not configured")
	}
	switch strings.ToLower(strings.TrimSpace(cfg.Provider)) {
	case "mock":
		return r.mock.StreamGenerate(ctx, prompt, cfg, history, onChunk)
	default:
		if cfg.APIKey == "" {
			log.Printf("llm: missing api key for model %s, falling back to mock", cfg.ID)
			return r.mock.StreamGenerate(ctx, prompt, cfg, history, onChunk)
		}
		return r.http.StreamGenerate(ctx, prompt, cfg, history, onChunk)
	}
}

// HTTPClient talks to OpenAI-compatible chat completion APIs.
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient wraps the provided http.Client or creates a default one.
func NewHTTPClient(httpClient *http.Client) *HTTPClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 150 * time.Second}
	}
	return &HTTPClient{client: httpClient}
}

func (c *HTTPClient) Generate(ctx context.Context, prompt string, cfg *model.ModelConfig, history []model.ChatMessage) (string, error) {
	if strings.TrimSpace(cfg.ModelName) == "" {
		return "", errors.New("model missing model_name")
	}
	baseURL := strings.TrimSpace(cfg.BaseURL)
	if baseURL == "" {
		baseURL = defaultAPIBase
	}
	baseURL = strings.TrimSuffix(baseURL, "/")
	messages := buildMessages(prompt, history)
	reqBody := map[string]interface{}{
		"model":    cfg.ModelName,
		"messages": messages,
	}
	if cfg.Temperature > 0 {
		reqBody["temperature"] = cfg.Temperature
	} else {
		reqBody["temperature"] = 0.8
	}
	if cfg.MaxTokens > 0 {
		reqBody["max_tokens"] = cfg.MaxTokens
	}
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		msg := fmt.Sprintf("llm provider error: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
		log.Printf("llm: %s", msg)
		return "", errors.New(msg)
	}
	var parsed completionResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		log.Printf("llm: decode error %v body=%s", err, string(body))
		return "", err
	}
	if parsed.Error != nil {
		return "", errors.New(parsed.Error.Message)
	}
	if len(parsed.Choices) == 0 {
		log.Printf("llm: empty choices body=%s", strings.TrimSpace(string(body)))
		return "", errors.New("llm returned no choices")
	}
	out := strings.TrimSpace(parsed.Choices[0].Message.Content)
	if out == "" {
		log.Printf("llm: empty content body=%s", strings.TrimSpace(string(body)))
	}
	return out, nil
}

func buildMessages(prompt string, history []model.ChatMessage) []completionMessage {
	messages := []completionMessage{}
	if strings.TrimSpace(prompt) != "" {
		messages = append(messages, completionMessage{Role: "system", Content: prompt})
	}
	if len(history) == 0 {
		// Some providers reject a chat with only system role; add a user turn to avoid “empty conversation”.
		messages = append(messages, completionMessage{
			Role:    "user",
			Content: "Generate based on the above instructions.",
		})
		return messages
	}
	for _, msg := range history {
		role := "user"
		if msg.Role == "assistant" {
			role = "assistant"
		} else if msg.Role == "system" {
			role = "system"
		}
		messages = append(messages, completionMessage{
			Role:    role,
			Content: msg.Content,
		})
	}
	return messages
}

type completionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type completionChoice struct {
	Message completionMessage `json:"message"`
}

type completionResponse struct {
	Choices []completionChoice `json:"choices"`
	Error   *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// StreamGenerate streams chunked responses (OpenAI-compatible stream).
func (c *HTTPClient) StreamGenerate(ctx context.Context, prompt string, cfg *model.ModelConfig, history []model.ChatMessage, onChunk func(contentDelta string, reasoningDelta string)) error {
	if strings.TrimSpace(cfg.ModelName) == "" {
		return errors.New("model missing model_name")
	}
	baseURL := strings.TrimSpace(cfg.BaseURL)
	if baseURL == "" {
		baseURL = defaultAPIBase
	}
	baseURL = strings.TrimSuffix(baseURL, "/")
	messages := buildMessages(prompt, history)
	reqBody := map[string]interface{}{
		"model":    cfg.ModelName,
		"messages": messages,
		"stream":   true,
	}
	if cfg.Temperature > 0 {
		reqBody["temperature"] = cfg.Temperature
	} else {
		reqBody["temperature"] = 0.8
	}
	if cfg.MaxTokens > 0 {
		reqBody["max_tokens"] = cfg.MaxTokens
	}
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		msg := fmt.Sprintf("llm provider error: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
		log.Printf("llm: %s", msg)
		return errors.New(msg)
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "data:") {
			continue
		}
		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if payload == "[DONE]" {
			break
		}
		var chunk struct {
			Choices []struct {
				Delta struct {
					Content   []string `json:"content,omitempty"`
					Reasoning []string `json:"reasoning_content,omitempty"`
				} `json:"delta"`
			} `json:"choices"`
			Error *struct {
				Message string `json:"message"`
			} `json:"error,omitempty"`
		}
		if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
			log.Printf("llm: stream decode err=%v payload=%s", err, payload)
			continue
		}
		if chunk.Error != nil {
			return errors.New(chunk.Error.Message)
		}
		if len(chunk.Choices) == 0 {
			continue
		}
		delta := chunk.Choices[0].Delta
		contentDelta := strings.Join(delta.Content, "")
		reasoningDelta := strings.Join(delta.Reasoning, "")
		if contentDelta != "" || reasoningDelta != "" {
			onChunk(contentDelta, reasoningDelta)
		}
	}
	return nil
}
