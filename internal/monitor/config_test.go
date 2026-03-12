package monitor

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// 1. Create a dummy config file with at least ONE target
	tmpFile := "test_targets.json"
	content := `{
        "interval_seconds": 30,
        "targets": [
            {
                "name": "Test-Target",
                "url": "http://localhost:8080",
                "expected_status": 200,
                "body_contains": "OK",
                "timeout_ms": 500,
                "retries": 1
            }
        ]
    }`

	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile) // Clean up after test

	// 2. Test loading the local temp file
	cfg, err := LoadConfig(tmpFile)
	if err != nil {
		t.Errorf("Expected successful load of temp config, got: %v", err)
	}

	// 3. Extra Check: Verify the data was actually parsed
	if len(cfg.Targets) == 0 {
		t.Error("Config loaded but targets array is empty")
	}
}
