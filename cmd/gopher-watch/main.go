package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/5x-Developer/gopher-watch-/internal/metrics"
	"github.com/5x-Developer/gopher-watch-/internal/monitor"
)

func main() {
	fmt.Println("Gopher-watch Monitoring Engine Starting")
	configPath := "configs/targets.json"
	// 1. Open (or create) the log file
	logFile, err := os.OpenFile("gopher-watch.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// 2. Create a JSON handler that writes to the file
	// We can also use io.MultiWriter(os.Stdout, logFile) to see it in both places
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger := slog.New(slog.NewJSONHandler(multiWriter, nil))

	logger.Info("Gopher-watch Monitoring Engine Starting", "version", "1.4")

	go func() {
		http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
			reg := metrics.GetInstance()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(reg)
		})
		fmt.Println("📈 Metrics endpoint available at http://localhost:8080/metrics")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// 1. Load config once to get the interval
	config, err := monitor.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("CRITICAL: %v", err)
	}

	ticker := time.NewTicker(time.Duration(config.IntervalSeconds) * time.Second)
	defer ticker.Stop()

	// 2. Setup Signal Handling for Graceful Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	runprobes := func() {
		// Re-load targets in case the file changed
		c, _ := monitor.LoadConfig(configPath)
		targets := c.Targets

		var wg sync.WaitGroup
		resultsChannel := make(chan monitor.Result, len(targets))

		fmt.Printf("\n--- Probe Cycle: %s ---\n", time.Now().Format("15:04:05"))

		for _, t := range targets {
			wg.Add(1)
			go func(target monitor.Target) {
				defer wg.Done()
				resultsChannel <- monitor.Ping(target)
			}(t)
		}

		go func() {
			wg.Wait()
			close(resultsChannel)
		}()

		total, passed := 0, 0
		for res := range resultsChannel {
			total++
			// 1. Update Global Metrics inside the loop
			metrics.GetInstance().RecordResult(res.TargetName, res.Success)

			if res.Success {
				passed++
				logger.Info("Probe Success",
					"target", res.TargetName,
					"status", res.Status,
					"latency_ms", res.Latency.Milliseconds(),
				)
			} else {
				logger.Error("Probe Failed",
					"target", res.TargetName,
					"status", res.Status,
					"latency_ms", res.Latency.Milliseconds(),
					"message", res.Message,
				)
			}
		}

		fmt.Printf("Summary: %d/%d passed\n", passed, total)
		metrics.GetInstance().UpdateTimestamp()
	}

	runprobes()

	for {
		select {
		case <-ticker.C:
			runprobes()
		case sig := <-sigChan:
			// 3. The Graceful Exit
			fmt.Printf("\nShutting down... Received signal: %v\n", sig)
			fmt.Println("Shutdown complete")
			return
		}
	}
}
