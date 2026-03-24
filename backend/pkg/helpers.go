package pkg

import "github.com/google/uuid"

// NewUUID generates a new UUID v4 string.
func NewUUID() string {
	return uuid.New().String()
}

// MaskSecret masks a secret string, showing only last 4 chars.
// e.g. "sk-ant-abc123xyz" → "sk-ant-****3xyz"
func MaskSecret(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	visible := s[len(s)-4:]
	return "****" + visible
}
