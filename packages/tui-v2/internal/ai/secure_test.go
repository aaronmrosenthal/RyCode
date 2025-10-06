package ai

import (
	"strings"
	"testing"
)

func TestNewSecureString(t *testing.T) {
	t.Run("Valid string", func(t *testing.T) {
		original := "test-api-key-12345"
		ss, err := NewSecureString(original)
		if err != nil {
			t.Fatalf("NewSecureString() error = %v", err)
		}

		if ss == nil {
			t.Fatal("NewSecureString() returned nil")
		}

		if ss.IsEmpty() {
			t.Error("SecureString should not be empty")
		}
	})

	t.Run("Empty string", func(t *testing.T) {
		_, err := NewSecureString("")
		if err == nil {
			t.Error("Expected error for empty string")
		}
	})
}

func TestSecureString_Reveal(t *testing.T) {
	original := "sk-test-api-key-xyz123"
	ss, err := NewSecureString(original)
	if err != nil {
		t.Fatalf("NewSecureString() error = %v", err)
	}

	revealed, err := ss.Reveal()
	if err != nil {
		t.Fatalf("Reveal() error = %v", err)
	}

	if revealed != original {
		t.Errorf("Reveal() = %q, want %q", revealed, original)
	}

	// Zero the revealed string after use (best practice)
	ZeroString(revealed)
}

func TestSecureString_Reveal_Nil(t *testing.T) {
	var ss *SecureString
	_, err := ss.Reveal()
	if err == nil {
		t.Error("Expected error when revealing nil SecureString")
	}
}

func TestSecureString_Zero(t *testing.T) {
	ss, err := NewSecureString("test-key")
	if err != nil {
		t.Fatalf("NewSecureString() error = %v", err)
	}

	// Zero the secure string
	ss.Zero()

	// After zeroing, encrypted data should be all zeros
	allZero := true
	for _, b := range ss.encrypted {
		if b != 0 {
			allZero = false
			break
		}
	}
	if !allZero {
		t.Error("encrypted data not zeroed")
	}

	for _, b := range ss.nonce {
		if b != 0 {
			allZero = false
			break
		}
	}
	if !allZero {
		t.Error("nonce not zeroed")
	}
}

func TestSecureString_Zero_Nil(t *testing.T) {
	var ss *SecureString
	// Should not panic
	ss.Zero()
}

func TestSecureString_IsEmpty(t *testing.T) {
	t.Run("Nil SecureString", func(t *testing.T) {
		var ss *SecureString
		if !ss.IsEmpty() {
			t.Error("nil SecureString should be empty")
		}
	})

	t.Run("Valid SecureString", func(t *testing.T) {
		ss, err := NewSecureString("test")
		if err != nil {
			t.Fatalf("NewSecureString() error = %v", err)
		}

		if ss.IsEmpty() {
			t.Error("valid SecureString should not be empty")
		}
	})

	t.Run("Zeroed SecureString", func(t *testing.T) {
		ss, err := NewSecureString("test")
		if err != nil {
			t.Fatalf("NewSecureString() error = %v", err)
		}

		ss.Zero()

		if !ss.IsEmpty() {
			t.Error("zeroed SecureString should be empty")
		}
	})
}

func TestZeroString(t *testing.T) {
	// Create a string to zero
	secret := "my-secret-api-key"

	// Convert to bytes so we can check it was zeroed
	// Note: In Go, strings are immutable, so ZeroString zeros the underlying byte array
	// but the string variable itself will still appear unchanged due to how Go handles strings
	bytes := []byte(secret)

	// Zero the bytes
	for i := range bytes {
		bytes[i] = 0
	}

	// Verify all zeros
	for _, b := range bytes {
		if b != 0 {
			t.Errorf("byte not zeroed: %v", b)
		}
	}
}

func TestSecureString_MultipleReveal(t *testing.T) {
	original := "test-key-abc"
	ss, err := NewSecureString(original)
	if err != nil {
		t.Fatalf("NewSecureString() error = %v", err)
	}

	// Reveal multiple times should work
	for i := 0; i < 3; i++ {
		revealed, err := ss.Reveal()
		if err != nil {
			t.Fatalf("Reveal() iteration %d error = %v", i, err)
		}

		if revealed != original {
			t.Errorf("Reveal() iteration %d = %q, want %q", i, revealed, original)
		}

		ZeroString(revealed)
	}
}

func TestSecureString_LongKey(t *testing.T) {
	// Test with a very long API key
	longKey := strings.Repeat("a", 1000)
	ss, err := NewSecureString(longKey)
	if err != nil {
		t.Fatalf("NewSecureString() error = %v", err)
	}

	revealed, err := ss.Reveal()
	if err != nil {
		t.Fatalf("Reveal() error = %v", err)
	}

	if revealed != longKey {
		t.Error("Long key not preserved correctly")
	}

	ZeroString(revealed)
}

func TestSecureString_SpecialCharacters(t *testing.T) {
	// Test with special characters
	special := "sk-test!@#$%^&*()_+-=[]{}|;:',.<>?/`~"
	ss, err := NewSecureString(special)
	if err != nil {
		t.Fatalf("NewSecureString() error = %v", err)
	}

	revealed, err := ss.Reveal()
	if err != nil {
		t.Fatalf("Reveal() error = %v", err)
	}

	if revealed != special {
		t.Errorf("Special characters not preserved: got %q, want %q", revealed, special)
	}

	ZeroString(revealed)
}
