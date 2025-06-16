package logger

import "testing"

func TestInitialize_InvalidLevel(t *testing.T) {
	err := Initialize("invalid-level")
	if err == nil {
		t.Fatal("expected error for invalid level, got nil")
	}
}
