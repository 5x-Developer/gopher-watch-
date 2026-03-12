package main

import (
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
	"github.com/5x-Developer/gopher-watch-/internal/notifier"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file into the system's environment
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}
}
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
			w.Header().Set("Content-Type", "text/plain; version=0.0.4")
			fmt.Fprint(w, metrics.GetInstance().ToPrometheus())
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

			// 1. Get the current streak before resetting it on success
			oldStreak := metrics.GetInstance().GetFailureStreak(res.TargetName)

			// 2. Record metrics and update/reset the streak
			metrics.GetInstance().RecordResult(res.TargetName, res.Success, res.Latency)
			streak := metrics.GetInstance().UpdateStreak(res.TargetName, res.Success)

			if res.Success {
				passed++
				if oldStreak >= 3 {
					recoveryMsg := fmt.Sprintf("✅ *Service Recovered*: %s\nEverything is back to normal after %d failed cycles.",
						res.TargetName, oldStreak)

					webhook := os.Getenv("SLACK_WEBHOOK_URL")
					if webhook != "" {
						go func(url, msg string) {
							if err := notifier.SendSlackAlert(url, msg); err != nil {
								logger.Error("Slack recovery alert failed", "error", err)
							}
						}(webhook, recoveryMsg)
					}
				}
				logger.Info("Probe Success",
					"target", res.TargetName,
					"status", res.Status,
					"latency_ms", res.Latency.Milliseconds(),
				)
			} else {
				logger.Error("Probe Failed",
					"target", res.TargetName,
					"status", res.Status,
					"streak", streak, // Added streak to logs for better debugging
					"message", res.Message,
				)

				// 3. SLACK ALERT LOGIC
				// Only alert exactly on the 3rd consecutive failed cycle (Tick) (9 attempts total)
				if streak == 3 {
					alertMsg := fmt.Sprintf("🚨 *Service Down*: %s\nFailed 3 consecutive cycles.\nError: %s",
						res.TargetName, res.Message)

					webhook := os.Getenv("SLACK_WEBHOOK_URL")
					if webhook != "" {
						go func(url, msg string) {
							if err := notifier.SendSlackAlert(url, msg); err != nil {
								logger.Error("Slack alert failed", "error", err)
							}
						}(webhook, alertMsg)
					}
				}
			}

		}
		fmt.Printf(" Summary: %d/%d passed\n", passed, total)
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
