package monitor

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Target struct {
	Name           string `json:"name"`
	URL            string `json:"url"`
	Method         string `json:"method"`
	ExpectedStatus int    `json:"expected_status"`
	BodyContains   string `json:"body_contains"`
	TimeoutMS      int    `json:"timeout_ms"`
}

func LoadTargetsFromFile(path string) ([]Target, error) {
	file, err := os.ReadFile(path)

	// throw: if file cannot be read
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// check: Is the file content just whitespace or empty?
	if len(strings.TrimSpace(string(file))) == 0 {
		return nil, fmt.Errorf("configuration file is empty")
	}

	var targets []Target
	/*Similar to a Arraylist in java, with a fixed type (Target)
	think of it as ArrayList<Target>, Target being our custom class that we defined*/

	err = json.Unmarshal(file, &targets)
	/*Similar to json.parse() in js
	Unmarshal takes a byte array
	the &targets is a pointer to the variable
	without the & it would be a copy of the variable, and any changes made to it would not affect the original variable*/

	//throw: if JSON is not in the expected format
	if err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	// Ensure we actually have at least one target to monitor
	if len(targets) == 0 {
		return nil, fmt.Errorf("no targets found in configuration")
	}

	return targets, nil
}
