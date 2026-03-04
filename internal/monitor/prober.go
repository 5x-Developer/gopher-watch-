package monitor

import (
	"fmt"
	"io"
	"math"
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
	var lastResult Result
	maxAttempts := t.Retries + 1

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		start := time.Now()

		client := http.Client{
			Timeout: time.Duration(t.TimeoutMS) * time.Millisecond,
		}

		response, err := client.Get(t.URL)
		latency := time.Since(start)

		var currentResult Result
		if err != nil {
			msg := err.Error()
			if os.IsTimeout(err) {
				msg = "Request timed out"
			}
			currentResult = Result{
				TargetName: t.Name, Success: false, Status: 0, Latency: latency,
				Message: fmt.Sprintf("Attempt %d/%d: %s", attempt, maxAttempts, msg),
			}
		} else {
			defer response.Body.Close()
			bodyBytes, _ := io.ReadAll(response.Body)
			bodyString := string(bodyBytes)

			statusOK := response.StatusCode == t.ExpectedStatus
			bodyOK := strings.Contains(bodyString, t.BodyContains)
			success := statusOK && bodyOK

			message := "OK"
			if !success {
				message = fmt.Sprintf("Attempt %d/%d: Fail (Status: %d, Body: %v)",
					attempt, maxAttempts, response.StatusCode, bodyOK)
			}

			currentResult = Result{
				TargetName: t.Name, Success: success, Status: response.StatusCode,
				Latency: latency, Message: message,
			}
		}

		lastResult = currentResult
		if currentResult.Success {
			return currentResult
		}

		// Exponential Backoff Logic
		if attempt < maxAttempts {
			waitTime := time.Duration(math.Pow(2, float64(attempt-1))) * time.Second
			time.Sleep(waitTime)
		}
	}
	return lastResult
}
