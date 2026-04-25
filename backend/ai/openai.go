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
	if strings.TrimSpace(baseURL) == "" {
		baseURL = defaultOpenAIBaseURL
	}
	return &OpenAIProvider{
		apiKey:    apiKey,
		model:     model,
		maxTokens: maxTokens,
		baseURL:   strings.TrimRight(baseURL, "/"),
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

func (o *OpenAIProvider) AnalyzeChat(ctx context.Context, systemPrompt string, chatTranscript string) (AIResponse, error) {
	return withRetry(ctx, "openai", func() (AIResponse, error) {
		reqBody := openAIResponsesRequest{
			Model:        o.model,
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

		model := parsed.Model
		if model == "" {
			model = o.model
		}

		return AIResponse{
			Content:      text,
			InputTokens:  parsed.Usage.InputTokens,
			OutputTokens: parsed.Usage.OutputTokens,
			Model:        model,
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
