package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	rps := flag.Int("rps", 10000, "requests per second")
	duration := flag.Duration("duration", 30*time.Second, "test duration")
	model := flag.String("model", "model1", "concurrency model: model1, model2, model3, model4, model5")
	workers := flag.Int("workers", 100, "number of workers (for fixed pool models)")
	// cpuWork := flag.Int("cpu", 1000, "CPU iterations per request")
	// ioWait := flag.Duration("io", 0, "simulated I/O wait per request")
	flag.Parse()

	log.Printf("HC-RE Starting")
	log.Printf("Model: %s | RPS: %d | Duration: %s | Workers: %d", *model, *rps, *duration, *workers)

	// I will do this later
	// Initialize system based on model
	// Run load test
	// Print metrics

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	fmt.Println("\nShutdown signal received")
}
