package ai

import "testing"

func TestPickCompatibleOpenAIModel(t *testing.T) {
	tests := []struct {
		name      string
		requested string
		available []string
		want      string
	}{
		{
			name:      "exact model id",
			requested: "gpt-5-mini",
			available: []string{
				"gpt-5-mini",
				"openai/gpt-5",
			},
			want: "gpt-5-mini",
		},
		{
			name:      "provider prefixed model id",
			requested: "gpt-5-mini",
			available: []string{
				"anthropic/claude-sonnet-4-6",
				"openai/gpt-5-mini",
			},
			want: "openai/gpt-5-mini",
		},
		{
			name:      "colon prefixed model id",
			requested: "gpt-5-mini",
			available: []string{
				"openai:gpt-5-mini",
			},
			want: "openai:gpt-5-mini",
		},
		{
			name:      "not available",
			requested: "gpt-5-mini",
			available: []string{
				"gpt-4.1-mini",
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PickCompatibleOpenAIModel(tt.requested, tt.available)
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
