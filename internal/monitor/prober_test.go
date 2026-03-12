package monitor

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing_Scenarios(t *testing.T) {
	// Create a local test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Gopher-Watch-Active"))
	}))
	defer ts.Close()

	tests := []struct {
		name           string
		expectedStatus int
		bodyContains   string
		expectSuccess  bool
	}{
		{"Successful Match", 200, "Gopher", true},
		{"Wrong Status Code", 404, "Gopher", false},
		{"Missing Body Text", 200, "Dragon", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := Target{
				URL:            ts.URL,
				ExpectedStatus: tt.expectedStatus,
				BodyContains:   tt.bodyContains,
				TimeoutMS:      500,
			}
			res := Ping(target)
			if res.Success != tt.expectSuccess {
				t.Errorf("%s: expected success %v, got %v", tt.name, tt.expectSuccess, res.Success)
			}
		})
	}
}
