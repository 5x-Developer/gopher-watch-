package monitor

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	IntervalSeconds int      `json:"interval_seconds"`
	Targets         []Target `json:"targets"`
}

type Target struct {
	IntervalSeconds int    `json:"interval_seconds"`
	Name            string `json:"name"`
	URL             string `json:"url"`
	Method          string `json:"method"`
	ExpectedStatus  int    `json:"expected_status"`
	BodyContains    string `json:"body_contains"`
	TimeoutMS       int    `json:"timeout_ms"`
	Retries         int    `json:"retries"`
}

// / LoadConfig reads the JSON configuration file and unmarshals it into a Config struct
func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)

	// throw: if file cannot be read
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// check: Is the file content just whitespace or empty?
	if len(strings.TrimSpace(string(file))) == 0 {
		return nil, fmt.Errorf("configuration file is empty")
	}

	var config Config
	/* Similar to an ArrayList in Java, with a fixed type (Config).
	   The Config struct now contains our slice of Targets.
	*/

	err = json.Unmarshal(file, &config)
	/* Similar to json.parse() in JS.
	   Unmarshal takes a byte array.
	   The &config is a pointer to the variable; without the &, it would be a
	   copy of the variable, and any changes made to it would not affect the original variable.
	*/

	// throw: if JSON is not in the expected format
	if err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	// Ensure we actually have at least one target to monitor
	if len(config.Targets) == 0 {
		return nil, fmt.Errorf("no targets found in configuration")
	}

	// SRE Best Practice: Ensure interval is safe (cannot be 0 or negative)
	if config.IntervalSeconds <= 0 {
		config.IntervalSeconds = 30 // Default fallback
	}

	return &config, nil
}
