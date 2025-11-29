package image

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/example/ai-avatar-studio/internal/model"
	llmclient "github.com/example/ai-avatar-studio/internal/pkg/llm"
	"github.com/example/ai-avatar-studio/internal/repository"
	"github.com/gorilla/websocket"
)

// Service orchestrates prompt generation and dispatch to image providers.
type Service struct {
	providers *repository.ImageProviderRepository
	presets   *repository.ImagePresetRepository
	jobs      *repository.ImageJobRepository
	chats     *repository.ChatRepository
	configs   *repository.ConfigRepository
	llm       llmclient.Client
	http      *http.Client
}

func NewService(
	providers *repository.ImageProviderRepository,
	presets *repository.ImagePresetRepository,
	jobs *repository.ImageJobRepository,
	chats *repository.ChatRepository,
	configs *repository.ConfigRepository,
	llm llmclient.Client,
) *Service {
	client := &http.Client{Timeout: 150 * time.Second}
	return &Service{
		providers: providers,
		presets:   presets,
		jobs:      jobs,
		chats:     chats,
		configs:   configs,
		llm:       llm,
		http:      client,
	}
}

// RequestImage triggers generation for a chat message/session.
func (s *Service) RequestImage(ctx context.Context, userID, sessionID, messageID, userPrompt string) (*model.ImageJob, error) {
	if s.providers == nil || s.presets == nil || s.jobs == nil || s.chats == nil {
		return nil, errors.New("image service unavailable")
	}
	provider, err := s.pickProvider(ctx)
	if err != nil {
		return nil, err
	}
	preset, _ := s.presets.Active(ctx)
	if preset == nil {
		return nil, errors.New("no active image preset")
	}

	finalPrompt, negativePrompt, err := s.buildPrompt(ctx, preset, sessionID, messageID, userPrompt)
	if err != nil {
		return nil, err
	}

	job := &model.ImageJob{
		UserID:         userID,
		SessionID:      sessionID,
		MessageID:      messageID,
		ProviderID:     provider.ID,
		PresetID:       preset.ID,
		Prompt:         userPrompt,
		NegativePrompt: negativePrompt,
		FinalPrompt:    finalPrompt,
		Status:         "pending",
	}
	job, err = s.jobs.Create(ctx, job)
	if err != nil {
		return nil, err
	}

	url, err := s.callProvider(ctx, provider, finalPrompt, negativePrompt)
	if err != nil {
		_ = s.jobs.UpdateStatus(ctx, job.ID, "failed", "", err.Error())
		job.Status = "failed"
		job.Error = err.Error()
		return job, err
	}
	job.Status = "succeeded"
	job.ResultURL = url
	job.UpdatedAt = time.Now()
	_ = s.jobs.UpdateStatus(ctx, job.ID, "succeeded", url, "")
	return job, nil
}

func (s *Service) GetJob(ctx context.Context, id string) (*model.ImageJob, error) {
	return s.jobs.Find(ctx, id)
}

func (s *Service) pickProvider(ctx context.Context) (*model.ImageProvider, error) {
	list, err := s.providers.ListActive(ctx)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("no active image provider")
	}
	if len(list) == 1 {
		return &list[0], nil
	}
	return &list[rand.Intn(len(list))], nil
}

type presetInstruction struct {
	Instruction string `json:"instruction"`
	Style       string `json:"style"`
	Negative    string `json:"negative"`
}

func (s *Service) buildPrompt(ctx context.Context, preset *model.ImagePreset, sessionID, messageID, userPrompt string) (string, string, error) {
	history, err := s.chats.ListMessages(ctx, sessionID, 20)
	if err != nil {
		return "", "", err
	}
	var convo strings.Builder
	for _, m := range history {
		convo.WriteString(fmt.Sprintf("%s: %s\n", m.Role, m.Content))
	}
	var presetCfg presetInstruction
	_ = json.Unmarshal([]byte(preset.PresetJSON), &presetCfg)
	system := strings.TrimSpace(presetCfg.Instruction)
	if system == "" {
		system = "You are an expert image prompt engineer. Summarize the dialogue into a concise English prompt for illustration."
	}
	style := strings.TrimSpace(presetCfg.Style)
	negative := strings.TrimSpace(presetCfg.Negative)

	builder := strings.Builder{}
	builder.WriteString(system)
	if style != "" {
		builder.WriteString("\nStyle hints: " + style)
	}
	builder.WriteString("\nConversation:\n")
	builder.WriteString(convo.String())
	if strings.TrimSpace(userPrompt) != "" {
		builder.WriteString("\nUser addition: " + strings.TrimSpace(userPrompt))
	}
	builder.WriteString("\nReturn only the final prompt for the image model.")

	modelCfg := s.resolvePromptModel(ctx, preset.PromptModelKey)
	if s.llm == nil {
		return "", "", errors.New("prompt llm unavailable")
	}
	reply, err := s.llm.Generate(ctx, builder.String(), modelCfg, nil)
	if err != nil {
		return "", "", fmt.Errorf("prompt llm error: %w", err)
	}
	final := strings.TrimSpace(reply)
	if final == "" {
		log.Printf("prompt llm empty reply (model=%s, session=%s, message=%s)", modelCfg.ID, sessionID, messageID)
		return "", "", errors.New("prompt llm returned empty")
	}
	return final, negative, nil
}

// --- helpers for NovelAI ---

func isNovelAIHost(base string) bool {
	l := strings.ToLower(strings.TrimSpace(base))
	return strings.Contains(l, "image.novelai.net") || strings.Contains(l, "api.novelai.net") || strings.Contains(l, "novelai.net")
}

func isV4FamilyModel(modelName string) bool {
	l := strings.ToLower(modelName)
	return strings.Contains(l, "nai-diffusion-4-5") ||
		strings.Contains(l, "nai-diffusion-4.5") ||
		strings.Contains(l, "nai-diffusion-4-curated") ||
		strings.Contains(l, "nai-diffusion-4-full") ||
		strings.Contains(l, "v4.5") ||
		strings.Contains(l, "v4 ")
}

func normalizeV4ModelName(modelName string) string {
	l := strings.ToLower(strings.TrimSpace(modelName))
	if strings.Contains(l, "4-5") || strings.Contains(l, "4.5") || strings.Contains(l, "v4.5") {
		if strings.Contains(l, "curated") {
			return "nai-diffusion-4-5-curated"
		}
		return "nai-diffusion-4-5-full"
	}
	if strings.Contains(l, "4-curated") || strings.Contains(l, "4_preview") || strings.Contains(l, "curated-preview") {
		return "nai-diffusion-4-curated-preview"
	}
	if strings.Contains(l, "4-full") || strings.Contains(l, "4 full") {
		return "nai-diffusion-4-full"
	}
	return modelName
}

func snap64(v int) int {
	if v <= 0 {
		return 64
	}
	v = (v / 64) * 64
	if v < 64 {
		v = 64
	}
	return v
}

func toInt(v interface{}, def int) int {
	switch t := v.(type) {
	case int:
		return t
	case int64:
		return int(t)
	case float64:
		return int(t)
	case json.Number:
		n, _ := t.Int64()
		return int(n)
	default:
		return def
	}
}

func toFloat(v interface{}, def float64) float64 {
	switch t := v.(type) {
	case float32:
		return float64(t)
	case float64:
		return t
	case int:
		return float64(t)
	case int64:
		return float64(t)
	case json.Number:
		f, _ := t.Float64()
		return f
	default:
		return def
	}
}

func toString(v interface{}, def string) string {
	if v == nil {
		return def
	}
	if s, ok := v.(string); ok {
		if strings.TrimSpace(s) == "" {
			return def
		}
		return s
	}
	return def
}

// callProvider routes to NAI v4 HTTP or other providers.
func (s *Service) callProvider(ctx context.Context, provider *model.ImageProvider, prompt, negative string) (string, error) {
	parameters := map[string]interface{}{}
	if strings.TrimSpace(provider.ParamsJSON) != "" {
		var params map[string]interface{}
		if err := json.Unmarshal([]byte(provider.ParamsJSON), &params); err == nil {
			for k, v := range params {
				parameters[k] = v
			}
		}
	}

	seedLocked := false
	if lock, ok := parameters["seed_lock"]; ok {
		if locked, ok := lock.(bool); ok && locked {
			seedLocked = true
		}
		delete(parameters, "seed_lock")
	}
	if !seedLocked {
		parameters["seed"] = time.Now().UnixNano() % 4294967295
	}

	modelName := strings.TrimSpace(provider.SelectedModel)
	if modelName == "" {
		modelName = "nai-diffusion-3"
	}

	baseURL := strings.TrimSpace(provider.BaseURL)
	lowerBase := strings.ToLower(baseURL)

	// NovelAI v4/v4.5: dedicated HTTP caller.
	if isNovelAIHost(baseURL) && isV4FamilyModel(modelName) {
		return s.callNovelAIV4HTTP(ctx, baseURL, provider.APIKey, modelName, prompt, negative, parameters)
	}

	// Non v4: keep negative prompt if provided.
	if negative != "" {
		parameters["negative_prompt"] = negative
	}

	// WebSocket provider (non NAI).
	if strings.HasPrefix(lowerBase, "ws") {
		return s.callProviderWS(ctx, baseURL, provider.APIKey, modelName, prompt, parameters)
	}

	// NovelAI v3 HTTP: auto append path if missing.
	if isNovelAIHost(baseURL) && !strings.Contains(lowerBase, "/ai/") {
		endpoint := strings.TrimRight(baseURL, "/") + "/ai/generate-image"
		return s.callProviderHTTP(ctx, endpoint, provider.APIKey, modelName, prompt, parameters)
	}

	// Other HTTP services: use as full endpoint.
	return s.callProviderHTTP(ctx, baseURL, provider.APIKey, modelName, prompt, parameters)
}

// NovelAI V4/V4.5 HTTP caller.
func (s *Service) callNovelAIV4HTTP(
	ctx context.Context,
	baseURL string,
	apiKey string,
	modelName string,
	prompt string,
	negative string,
	parameters map[string]interface{},
) (string, error) {
	base := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if !strings.Contains(strings.ToLower(base), "/ai/") {
		base = base + "/ai/generate-image"
	}

	trimmedPrompt := strings.TrimSpace(prompt)
	if strings.Contains(trimmedPrompt, "|||") {
		parts := strings.Split(trimmedPrompt, "|||")
		if len(parts) > 0 {
			trimmedPrompt = strings.TrimSpace(parts[0])
		}
	}
	trimmedNegative := strings.TrimSpace(negative)

	width := snap64(toInt(parameters["width"], 832))
	height := snap64(toInt(parameters["height"], 1216))
	scale := toFloat(parameters["scale"], 5.0)
	if scale > 10 {
		scale = 10
	}
	// For V4.5 stick to the known-good sampler/steps from manual success.
	steps := 23
	sampler := "k_euler_ancestral"
	noiseSchedule := toString(parameters["noise_schedule"], "karras")
	ucPreset := toInt(parameters["ucPreset"], 0)
	seed := toInt(parameters["seed"], int(time.Now().UnixNano()%4294967295))

	// Build payload matching the proven working shape.
	params := map[string]interface{}{
		"params_version":      3,
		"width":               width,
		"height":              height,
		"scale":               scale,
		"sampler":             sampler,
		"steps":               steps,
		"n_samples":           1,
		"ucPreset":            ucPreset,
		"qualityToggle":       true,
		"dynamic_thresholding": false,
		"controlnet_strength":  1,
		"legacy":               false,
		"add_original_image":   true,
		"cfg_rescale":          0,
		"noise_schedule":       noiseSchedule,
		"legacy_v3_extend":     false,
		"skip_cfg_above_sigma": nil,
		"use_coords":           false,
		"seed":                 seed,
		"characterPrompts":     []interface{}{},
		"reference_image_multiple":             []interface{}{},
		"reference_information_extracted_multiple": []interface{}{},
		"reference_strength_multiple":          []interface{}{},
		"deliberate_euler_ancestral_bug":       false,
		"prefer_brownian":                      true,
		"negative_prompt": func() string {
			if trimmedNegative == "" {
				return ""
			}
			return trimmedNegative
		}(),
		"v4_prompt": map[string]interface{}{
			"caption": map[string]interface{}{
				"base_caption":  trimmedPrompt,
				"char_captions": []interface{}{},
			},
			"use_coords": false,
			"use_order":  true,
		},
		"v4_negative_prompt": map[string]interface{}{
			"caption": map[string]interface{}{
				"base_caption":  trimmedNegative,
				"char_captions": []interface{}{},
			},
		},
	}

	body := map[string]interface{}{
		"action":     "generate",
		"input":      trimmedPrompt,
		"model":      normalizeV4ModelName(modelName),
		"parameters": params,
	}

	payload, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "*/*")

	resp, err := s.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		log.Printf("NovelAI V4 image provider error model=%s status=%d body=%s payload_params=%v",
			modelName, resp.StatusCode, strings.TrimSpace(string(bodyBytes)), params)
		return "", fmt.Errorf("provider error status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return extractImageFromZip(data)
}

func (s *Service) resolvePromptModel(ctx context.Context, modelKey string) *model.ModelConfig {
	// Use a real model with API key when available, otherwise fall back to mock.
	if s.configs == nil {
		return &model.ModelConfig{ID: "prompt-generator", Provider: "mock", Status: "active"}
	}
	key := strings.TrimSpace(modelKey)
	if key != "" {
		if m, err := s.configs.FindModel(ctx, key); err == nil && m != nil {
			return m
		}
	}
	if m, err := s.configs.DefaultModel(ctx); err == nil && m != nil {
		return m
	}
	return &model.ModelConfig{ID: "prompt-generator", Provider: "mock", Status: "active"}
}

func toWSBase(raw string) string {
	u := strings.TrimSpace(raw)
	l := strings.ToLower(u)
	if strings.HasPrefix(l, "http://") {
		return "ws://" + strings.TrimPrefix(u, "http://")
	}
	if strings.HasPrefix(l, "https://") {
		return "wss://" + strings.TrimPrefix(u, "https://")
	}
	if strings.HasPrefix(l, "ws") {
		return u
	}
	return "wss://" + strings.TrimPrefix(u, "//")
}

// Generic HTTP provider (non-NAI v4).
func (s *Service) callProviderHTTP(ctx context.Context, baseURL, apiKey, modelName, prompt string, parameters map[string]interface{}) (string, error) {
	body := map[string]interface{}{
		"action":     "generate",
		"input":      prompt,
		"model":      modelName,
		"parameters": parameters,
	}
	payload, _ := json.Marshal(body)
	endpoint := strings.TrimSpace(baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if strings.TrimSpace(apiKey) != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	req.Header.Set("Accept", "application/zip")
	resp, err := s.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		log.Printf("image provider error model=%s status=%d body=%s payload_params=%v",
			modelName, resp.StatusCode, strings.TrimSpace(string(bodyBytes)), parameters)
		return "", fmt.Errorf("provider error status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return extractImageFromZip(data)
}

// WebSocket provider (non-NAI).
func (s *Service) callProviderWS(ctx context.Context, baseURL, apiKey, modelName, prompt string, parameters map[string]interface{}) (string, error) {
	endpoint := strings.TrimSpace(baseURL)
	if endpoint == "" {
		return "", errors.New("provider base url is empty")
	}
	dialer := websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		EnableCompression: true,
	}
	header := http.Header{}
	if strings.TrimSpace(apiKey) != "" {
		header.Set("Authorization", "Bearer "+apiKey)
	}
	conn, _, err := dialer.DialContext(ctx, endpoint, header)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	payload := map[string]interface{}{
		"action":     "generateImage",
		"model":      modelName,
		"input":      prompt,
		"parameters": parameters,
	}
	msg, _ := json.Marshal(payload)
	if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		return "", err
	}

	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			return "", err
		}
		if msgType == websocket.BinaryMessage {
			return extractImageFromZip(data)
		}
		if msgType == websocket.TextMessage {
			var parsed map[string]interface{}
			if err := json.Unmarshal(data, &parsed); err == nil {
				if msg, ok := parsed["message"].(string); ok {
					return "", fmt.Errorf("provider error: %s", msg)
				}
			}
			return "", fmt.Errorf("provider error: %s", strings.TrimSpace(string(data)))
		}
	}
}

func extractImageFromZip(data []byte) (string, error) {
	readerAt := bytes.NewReader(data)
	zr, err := zip.NewReader(readerAt, int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("decode zip: %w", err)
	}
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			continue
		}
		imgBytes, err := io.ReadAll(rc)
		rc.Close()
		if err != nil || len(imgBytes) == 0 {
			continue
		}
		encoded := base64.StdEncoding.EncodeToString(imgBytes)
		return "data:image/png;base64," + encoded, nil
	}
	return "", errors.New("no image in provider response")
}
