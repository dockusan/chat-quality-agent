package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const defaultOpenAIBaseURL = "https://api.openai.com/v1"

type OpenAIProvider struct {
	apiKey    string
	model     string
	maxTokens int
	baseURL   string
}

func NewOpenAIProvider(apiKey, model string, maxTokens int, baseURL string) *OpenAIProvider {
	if model == "" {
		model = "gpt-5-mini"
	}
	if maxTokens <= 0 {
		maxTokens = 16384
	}
	return &OpenAIProvider{
		apiKey:    apiKey,
		model:     model,
		maxTokens: maxTokens,
		baseURL:   normalizeOpenAIBaseURL(baseURL),
	}
}

type openAIResponsesRequest struct {
	Model           string        `json:"model"`
	Instructions    string        `json:"instructions,omitempty"`
	Input           []openAIInput `json:"input"`
	MaxOutputTokens int           `json:"max_output_tokens,omitempty"`
	Text            openAIText    `json:"text"`
}

type openAIText struct {
	Format openAITextFormat `json:"format"`
}

type openAITextFormat struct {
	Type string `json:"type"`
}

type openAIInput struct {
	Role    string               `json:"role"`
	Content []openAIInputContent `json:"content"`
}

type openAIInputContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type openAIResponsesResponse struct {
	Model      string               `json:"model"`
	OutputText string               `json:"output_text"`
	Output     []openAIOutputItem   `json:"output"`
	Usage      openAIUsage          `json:"usage"`
	Error      *openAIErrorResponse `json:"error"`
}

type openAIOutputItem struct {
	Type    string                `json:"type"`
	Content []openAIOutputContent `json:"content"`
}

type openAIOutputContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type openAIUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type openAIErrorResponse struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

type openAIModelsResponse struct {
	Data []openAIModel `json:"data"`
}

type openAIModel struct {
	ID string `json:"id"`
}

func (o *OpenAIProvider) AnalyzeChat(ctx context.Context, systemPrompt string, chatTranscript string) (AIResponse, error) {
	return withRetry(ctx, "openai", func() (AIResponse, error) {
		model, err := o.resolveModel(ctx)
		if err != nil {
			return AIResponse{}, err
		}

		reqBody := openAIResponsesRequest{
			Model:        model,
			Instructions: systemPrompt,
			Input: []openAIInput{
				{
					Role: "user",
					Content: []openAIInputContent{
						{
							Type: "input_text",
							Text: chatTranscript,
						},
					},
				},
			},
			MaxOutputTokens: o.maxTokens,
			Text: openAIText{
				Format: openAITextFormat{Type: "text"},
			},
		}

		payload, err := json.Marshal(reqBody)
		if err != nil {
			return AIResponse{}, fmt.Errorf("openai request marshal error: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.baseURL+"/responses", bytes.NewReader(payload))
		if err != nil {
			return AIResponse{}, fmt.Errorf("openai request build error: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+o.apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := NewHTTPClientWithTimeout().Do(req)
		if err != nil {
			return AIResponse{}, fmt.Errorf("openai api error: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return AIResponse{}, fmt.Errorf("openai response read error: %w", err)
		}

		var parsed openAIResponsesResponse
		if err := json.Unmarshal(body, &parsed); err != nil {
			return AIResponse{}, fmt.Errorf("openai response parse error: %w", err)
		}

		if resp.StatusCode >= http.StatusBadRequest {
			if parsed.Error != nil && parsed.Error.Message != "" {
				return AIResponse{}, fmt.Errorf("openai api error: %s", parsed.Error.Message)
			}
			return AIResponse{}, fmt.Errorf("openai api error: status %d", resp.StatusCode)
		}

		text := strings.TrimSpace(parsed.OutputText)
		if text == "" {
			text = extractOpenAIOutputText(parsed.Output)
		}
		if text == "" {
			return AIResponse{}, fmt.Errorf("openai api returned empty content")
		}

		responseModel := parsed.Model
		if responseModel == "" {
			responseModel = reqBody.Model
		}

		return AIResponse{
			Content:      text,
			InputTokens:  parsed.Usage.InputTokens,
			OutputTokens: parsed.Usage.OutputTokens,
			Model:        responseModel,
			Provider:     "openai",
		}, nil
	})
}

func (o *OpenAIProvider) AnalyzeChatBatch(ctx context.Context, systemPrompt string, items []BatchItem) (AIResponse, error) {
	batchPrompt := WrapBatchPrompt(systemPrompt, len(items))
	batchTranscript := FormatBatchTranscript(items)
	return o.AnalyzeChat(ctx, batchPrompt, batchTranscript)
}

func extractOpenAIOutputText(items []openAIOutputItem) string {
	var parts []string
	for _, item := range items {
		if item.Type != "message" {
			continue
		}
		for _, content := range item.Content {
			if content.Type == "output_text" && strings.TrimSpace(content.Text) != "" {
				parts = append(parts, strings.TrimSpace(content.Text))
			}
		}
	}
	return strings.Join(parts, "\n")
}

func (o *OpenAIProvider) resolveModel(ctx context.Context) (string, error) {
	if o.baseURL == defaultOpenAIBaseURL {
		return o.model, nil
	}

	models, err := FetchOpenAIModels(ctx, o.apiKey, o.baseURL)
	if err != nil {
		return "", err
	}
	if len(models) == 0 {
		return o.model, nil
	}

	if model := PickCompatibleOpenAIModel(o.model, models); model != "" {
		return model, nil
	}

	return "", fmt.Errorf("openai model %q is not available from %s; available examples: %s", o.model, o.baseURL, strings.Join(limitStrings(models, 8), ", "))
}

// FetchOpenAIModels reads model IDs from an OpenAI-compatible /models endpoint.
func FetchOpenAIModels(ctx context.Context, apiKey, baseURL string) ([]string, error) {
	baseURL = normalizeOpenAIBaseURL(baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/models", nil)
	if err != nil {
		return nil, fmt.Errorf("openai models request build error: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := NewHTTPClientWithTimeout().Do(req)
	if err != nil {
		return nil, fmt.Errorf("openai models api error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("openai models response read error: %w", err)
	}

	var parsed openAIModelsResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("openai models response parse error: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		var errorResp struct {
			Error *openAIErrorResponse `json:"error"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Error != nil && errorResp.Error.Message != "" {
			return nil, fmt.Errorf("openai models api error: %s", errorResp.Error.Message)
		}
		return nil, fmt.Errorf("openai models api error: status %d", resp.StatusCode)
	}

	models := make([]string, 0, len(parsed.Data))
	for _, model := range parsed.Data {
		if strings.TrimSpace(model.ID) != "" {
			models = append(models, model.ID)
		}
	}
	return models, nil
}

// PickCompatibleOpenAIModel keeps exact model IDs first, then supports proxies
// that expose provider-prefixed IDs such as "openai/gpt-5-mini".
func PickCompatibleOpenAIModel(requested string, available []string) string {
	requested = strings.TrimSpace(requested)
	if requested == "" {
		return ""
	}

	for _, model := range available {
		if model == requested {
			return model
		}
	}
	for _, model := range available {
		if strings.HasSuffix(model, "/"+requested) || strings.HasSuffix(model, ":"+requested) {
			return model
		}
	}
	return ""
}

func normalizeOpenAIBaseURL(baseURL string) string {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		return defaultOpenAIBaseURL
	}
	return baseURL
}

func limitStrings(values []string, limit int) []string {
	if len(values) <= limit {
		return values
	}
	return values[:limit]
}
