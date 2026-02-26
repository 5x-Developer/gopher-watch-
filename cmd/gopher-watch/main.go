package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/5x-Developer/gopher-watch-/internal/monitor"
)

func main() {
	fmt.Println("Gopher-watch Monitoring Engine Starting")

	configPath := "configs/targets.json"
	targets, err := monitor.LoadTargetsFromFile(configPath)
	if err != nil {
		// Log the error and exit with code 1 so monitoring/CI-CD knows it failed
		log.Fatalf("CRITICAL: Application failed to initialize: %v", err)
	}

	var wg sync.WaitGroup                                     // WaitGroup to wait for all probes to finish
	resultsChannel := make(chan monitor.Result, len(targets)) // Buffered channel to collect results
	fmt.Printf("Successfully loaded %d targets:\n", len(targets))
	for _, t := range targets {
		wg.Add(1)
		go func(target monitor.Target) {
			defer wg.Done()
			//sending the result of the ping to the results channel, which will be collected by the main goroutine
			resultsChannel <- monitor.Ping(target)
		}(t)
	}
	go func() {
		wg.Wait()             // Wait for all probes to finish
		close(resultsChannel) // Close the channel to signal that no more results will be sent
	}()
	// 4. The "Collector": Pull results from the channel and print them
	// This loop stays alive until the channel is closed
	total := 0
	passed := 0
	failed := 0
	for res := range resultsChannel {
		total++
		if res.Success {
			passed++
			fmt.Printf("%s: %d | Latency: %v\n", res.TargetName, res.Status, res.Latency.Round(time.Millisecond))
		} else {
			failed++
			fmt.Printf("%s: FAILED | Error: %s\n", res.TargetName, res.Message)
		}
	}
	fmt.Println("All probes completed. Exiting.")
	fmt.Println("--------------------------------------------------")
	fmt.Printf(" Monitoring Summary:\n")
	fmt.Printf("   Total Probes:  %d\n", total)
	fmt.Printf("   Passed:        %d\n", passed)
	fmt.Printf("   Failed:        %d\n", failed)
	fmt.Println("--------------------------------------------------")
}
