package notifications

import (
	"testing"
)

func TestNewTelegramNotifier(t *testing.T) {
	notifier := NewTelegramNotifier("test-bot-token", "-1001234567890")
	if notifier == nil {
		t.Fatal("Notifier should not be nil")
	}
	if notifier.botToken != "test-bot-token" {
		t.Errorf("Expected bot token 'test-bot-token', got %s", notifier.botToken)
	}
	if notifier.chatID != "-1001234567890" {
		t.Errorf("Expected chat ID '-1001234567890', got %s", notifier.chatID)
	}
}

func TestNewEmailNotifier(t *testing.T) {
	notifier := NewEmailNotifier("smtp.gmail.com", 587, "user", "pass", "from@test.com", []string{"to@test.com"})
	if notifier == nil {
		t.Fatal("Notifier should not be nil")
	}
	if notifier.smtpHost != "smtp.gmail.com" {
		t.Errorf("Expected smtp host 'smtp.gmail.com', got %s", notifier.smtpHost)
	}
}

func TestSplitComma(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"a@b.com, c@d.com", 2},
		{"single@email.com", 1},
		{"a@b.com,c@d.com, e@f.com", 3},
		{"", 0}, // empty string returns 0 items after trim
	}

	for _, tt := range tests {
		result := splitComma(tt.input)
		if len(result) != tt.expected {
			t.Errorf("splitComma(%q) = %d items, want %d", tt.input, len(result), tt.expected)
		}
	}
}
