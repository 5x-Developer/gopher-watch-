package metrics

import (
	"fmt"
	"sync"
	"time"
)

// Registry holds the aggregated stats for the entire application life
type Registry struct {
	mu                  sync.RWMutex
	TotalProbes         int            `json:"total_probes"`
	TotalPassed         int            `json:"total_passed"`
	TotalFailed         int            `json:"total_failed"`
	TotalLatency        int64          `json:"total_latency_ms"` // Cumulative latency for averages
	LastProbeAt         string         `json:"last_probe_at"`
	TargetStats         map[string]int `json:"target_failures"`
	ConsecutiveFailures map[string]int `json:"consecutive_failures"` // Track failure streaks
}

var (
	instance *Registry
	once     sync.Once
)

// GetInstance returns a singleton of the metrics registry
func GetInstance() *Registry {
	once.Do(func() {
		instance = &Registry{
			TargetStats:         make(map[string]int),
			ConsecutiveFailures: make(map[string]int),
		}
	})
	return instance
}

// RecordResult updates the counts for each probe
func (r *Registry) RecordResult(target string, success bool, latency time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.TotalProbes++
	if success {
		r.TotalPassed++
		r.TotalLatency += latency.Milliseconds() // Add to total
	} else {
		r.TotalFailed++
		r.TargetStats[target]++
	}
}

// UpdateTimestamp sets the last time a cycle was completed
func (r *Registry) UpdateTimestamp() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.LastProbeAt = time.Now().Format(time.RFC3339)
}

// Format the metrics in Prometheus exposition format
func (r *Registry) ToPrometheus() string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return fmt.Sprintf(
		"# HELP gopher_watch_total_probes Total probes executed\n"+
			"# TYPE gopher_watch_total_probes counter\n"+
			"gopher_watch_total_probes %d\n\n"+
			"# HELP gopher_watch_total_passed Total passed probes\n"+
			"# TYPE gopher_watch_total_passed counter\n"+
			"gopher_watch_total_passed %d\n\n"+
			"# HELP gopher_watch_total_failed Total failed probes\n"+
			"# TYPE gopher_watch_total_failed counter\n"+
			"gopher_watch_total_failed %d\n\n"+
			"# HELP gopher_watch_latency_ms Total latency in ms\n"+
			"# TYPE gopher_watch_latency_ms counter\n"+
			"gopher_watch_latency_ms %d\n",
		r.TotalProbes, r.TotalPassed, r.TotalFailed, r.TotalLatency,
	)
}

func (r *Registry) GetFailureStreak(target string) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.ConsecutiveFailures[target]
}

// UpdateStreak increments on failure and resets on success. Returns the current streak.
func (r *Registry) UpdateStreak(target string, success bool) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	if success {
		r.ConsecutiveFailures[target] = 0
		return 0
	}

	r.ConsecutiveFailures[target]++
	return r.ConsecutiveFailures[target]
}
