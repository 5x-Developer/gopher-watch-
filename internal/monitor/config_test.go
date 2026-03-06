package monitor

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// 1. Test missing file
	_, err := LoadConfig("non-existent.json")
	if err == nil {
		t.Error("Expected error for missing file, got nil")
	}

	// 2. Test invalid JSON
	tmpFile := "test_bad.json"
	os.WriteFile(tmpFile, []byte("{ bad json "), 0644)
	defer os.Remove(tmpFile)

	_, err = LoadConfig(tmpFile)
	if err == nil {
		t.Error("Expected error for malformed JSON, got nil")
	}
}
