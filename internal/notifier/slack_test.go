package notifier

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendSlackAlert(t *testing.T) {
	// Mock Slack API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	err := SendSlackAlert(server.URL, "Test Message")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
