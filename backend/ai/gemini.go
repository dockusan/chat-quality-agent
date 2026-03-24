package ai

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

type GeminiProvider struct {
	apiKey string
	model  string
}

func NewGeminiProvider(apiKey, model string) *GeminiProvider {
	if model == "" {
		model = "gemini-2.0-flash"
	}
	return &GeminiProvider{
		apiKey: apiKey,
		model:  model,
	}
}

func (g *GeminiProvider) AnalyzeChat(ctx context.Context, systemPrompt string, chatTranscript string) (AIResponse, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  g.apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return AIResponse{}, fmt.Errorf("gemini client error: %w", err)
	}

	result, err := client.Models.GenerateContent(ctx, g.model, genai.Text(chatTranscript), &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(systemPrompt, "user"),
	})
	if err != nil {
		return AIResponse{}, fmt.Errorf("gemini api error: %w", err)
	}

	text := result.Text()
	if text == "" {
		return AIResponse{}, fmt.Errorf("gemini api returned empty content")
	}

	aiResp := AIResponse{
		Content:  text,
		Model:    g.model,
		Provider: "gemini",
	}
	if result.UsageMetadata != nil {
		aiResp.InputTokens = int(result.UsageMetadata.PromptTokenCount)
		aiResp.OutputTokens = int(result.UsageMetadata.CandidatesTokenCount)
	}

	return aiResp, nil
}

func (g *GeminiProvider) AnalyzeChatBatch(ctx context.Context, systemPrompt string, items []BatchItem) (AIResponse, error) {
	batchPrompt := WrapBatchPrompt(systemPrompt, len(items))
	batchTranscript := FormatBatchTranscript(items)
	return g.AnalyzeChat(ctx, batchPrompt, batchTranscript)
}
