package metrics

import (
	"sync"
	"time"
)

// Registry holds the aggregated stats for the entire application life
type Registry struct {
	mu          sync.RWMutex
	TotalProbes int            `json:"total_probes"`
	TotalPassed int            `json:"total_passed"`
	TotalFailed int            `json:"total_failed"`
	LastProbeAt string         `json:"last_probe_at"`
	TargetStats map[string]int `json:"target_failures"`
}

var (
	instance *Registry
	once     sync.Once
)

// GetInstance returns a singleton of the metrics registry
func GetInstance() *Registry {
	once.Do(func() {
		instance = &Registry{
			TargetStats: make(map[string]int),
		}
	})
	return instance
}

// RecordResult updates the counts for each probe
func (r *Registry) RecordResult(target string, success bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.TotalProbes++
	if success {
		r.TotalPassed++
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
