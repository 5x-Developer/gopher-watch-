package monitor

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Result struct {
	TargetName string        `json:"target_name"`
	Success    bool          `json:"success"`
	Status     int           `json:"status"`
	Latency    time.Duration `json:"latency"`
	Message    string        `json:"message,omitempty"`
}

func Ping(t Target) Result {
	start := time.Now()

	client := http.Client{
		Timeout: time.Duration(t.TimeoutMS) * time.Millisecond,
	}
	response, err := client.Get(t.URL)
	latency := time.Since(start)

	if err != nil {
		msg := err.Error()
		if os.IsTimeout(err) {
			msg = "Request timed out"
		}
		return Result{
			TargetName: t.Name,
			Success:    false,
			Status:     0,
			Latency:    latency,
			Message:    msg,
		}
	}
	defer response.Body.Close()

	// 1. Read the response body into memory
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return Result{
			TargetName: t.Name,
			Success:    false,
			Status:     response.StatusCode,
			Latency:    latency,
			Message:    "Failed to read response body: " + err.Error(),
		}
	}

	bodyString := string(bodyBytes)

	// 2. Perform the checks
	statusOK := response.StatusCode == t.ExpectedStatus
	bodyOK := strings.Contains(bodyString, t.BodyContains)

	// 3. Determine overall success
	success := statusOK && bodyOK

	// 4. Construct a descriptive message if it fails
	message := "OK"
	if !statusOK {
		message = fmt.Sprintf("Status mismatch: expected %d, got %d", t.ExpectedStatus, response.StatusCode)
	} else if !bodyOK {
		message = fmt.Sprintf("Body assertion failed: could not find '%s'", t.BodyContains)
	}

	return Result{
		TargetName: t.Name,
		Success:    success,
		Status:     response.StatusCode,
		Latency:    latency,
		Message:    message,
	}
}
