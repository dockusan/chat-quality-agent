package pkg

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key := "test-encryption-key-32-bytes!!" // Must be 32 bytes for AES-256
	// Pad to 32 bytes
	for len(key) < 32 {
		key += "x"
	}
	key = key[:32]

	plaintext := []byte("Hello, secret world!")

	encrypted, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if string(encrypted) == string(plaintext) {
		t.Fatal("Encrypted data should differ from plaintext")
	}

	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Fatalf("Decrypted data mismatch: got %q, want %q", decrypted, plaintext)
	}
}

func TestEncryptDecryptBase64(t *testing.T) {
	key := "12345678901234567890123456789012" // Exactly 32 bytes

	plaintext := []byte("API key: sk-ant-abc123")

	encoded, err := EncryptToBase64(plaintext, key)
	if err != nil {
		t.Fatalf("EncryptToBase64 failed: %v", err)
	}

	decoded, err := DecryptFromBase64(encoded, key)
	if err != nil {
		t.Fatalf("DecryptFromBase64 failed: %v", err)
	}

	if string(decoded) != string(plaintext) {
		t.Fatalf("Round-trip mismatch: got %q, want %q", decoded, plaintext)
	}
}

func TestDecryptWrongKey(t *testing.T) {
	key1 := "12345678901234567890123456789012"
	key2 := "abcdefghijklmnopqrstuvwxyz123456"

	plaintext := []byte("secret data")
	encrypted, _ := Encrypt(plaintext, key1)

	_, err := Decrypt(encrypted, key2)
	if err == nil {
		t.Fatal("Decrypt with wrong key should fail")
	}
}

func TestNewUUID(t *testing.T) {
	id1 := NewUUID()
	id2 := NewUUID()

	if id1 == "" || id2 == "" {
		t.Fatal("UUID should not be empty")
	}
	if id1 == id2 {
		t.Fatal("Two UUIDs should be different")
	}
	if len(id1) != 36 {
		t.Fatalf("UUID length should be 36, got %d", len(id1))
	}
}

func TestMaskSecret(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"sk-ant-api03-abc123xyz", "****3xyz"},
		{"short", "****hort"},
		{"ab", "****"},
		{"", "****"},
	}

	for _, tt := range tests {
		result := MaskSecret(tt.input)
		if result != tt.expected {
			t.Errorf("MaskSecret(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
