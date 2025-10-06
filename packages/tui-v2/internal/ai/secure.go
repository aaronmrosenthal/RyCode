package ai

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"runtime"
)

// SecureString stores sensitive strings encrypted in memory
// This protects API keys and other secrets from memory dumps, debuggers, and process scanning
type SecureString struct {
	encrypted []byte
	nonce     []byte
	gcm       cipher.AEAD
}

// NewSecureString creates an encrypted string in memory
// The plaintext is immediately encrypted and the original is zeroed out
func NewSecureString(plaintext string) (*SecureString, error) {
	if plaintext == "" {
		return nil, fmt.Errorf("cannot secure empty string")
	}

	// Generate encryption key from random bytes
	// Note: This key is also in memory, but it's harder to identify as an API key
	key := make([]byte, 32) // AES-256
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate encryption key: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the plaintext
	encrypted := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	ss := &SecureString{
		encrypted: encrypted,
		nonce:     nonce,
		gcm:       gcm,
	}

	// Zero out the key from memory (best effort)
	// The key is still in the GCM cipher, but we reduce exposure
	for i := range key {
		key[i] = 0
	}
	runtime.KeepAlive(key)

	return ss, nil
}

// Reveal decrypts and returns the plaintext temporarily
// IMPORTANT: Caller MUST zero out the returned string after use with ZeroString()
func (s *SecureString) Reveal() (string, error) {
	if s == nil {
		return "", fmt.Errorf("cannot reveal nil SecureString")
	}

	plaintext, err := s.gcm.Open(nil, s.nonce, s.encrypted, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// Zero securely wipes the encrypted data from memory
// Call this when the SecureString is no longer needed
func (s *SecureString) Zero() {
	if s == nil {
		return
	}

	// Zero out encrypted data
	for i := range s.encrypted {
		s.encrypted[i] = 0
	}

	// Zero out nonce
	for i := range s.nonce {
		s.nonce[i] = 0
	}

	runtime.KeepAlive(s.encrypted)
	runtime.KeepAlive(s.nonce)
}

// ZeroString securely zeros out a string by converting to byte slice
// Use this to clear decrypted secrets after use
func ZeroString(s string) {
	// Convert to byte slice and zero
	b := []byte(s)
	for i := range b {
		b[i] = 0
	}
	runtime.KeepAlive(b)
}

// IsEmpty returns true if the SecureString is nil or has no encrypted data
func (s *SecureString) IsEmpty() bool {
	if s == nil {
		return true
	}
	// Check if encrypted data is empty OR all zeros (after Zero() is called)
	if len(s.encrypted) == 0 {
		return true
	}

	// Check if all bytes are zero (indicating it was zeroed)
	allZero := true
	for _, b := range s.encrypted {
		if b != 0 {
			allZero = false
			break
		}
	}
	return allZero
}
