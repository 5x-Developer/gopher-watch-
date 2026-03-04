package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/5x-Developer/gopher-watch-/internal/monitor"
)

func main() {
	fmt.Println("Gopher-watch Monitoring Engine Starting")
	configPath := "configs/targets.json"

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
			if res.Success {
				passed++
				fmt.Printf("%s: %d | %v\n", res.TargetName, res.Status, res.Latency.Round(time.Millisecond))
			} else {
				fmt.Printf("%s: FAILED | %s\n", res.TargetName, res.Message)
			}
		}
		fmt.Printf("Summary: %d/%d passed\n", passed, total)
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
