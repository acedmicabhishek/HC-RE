package pkg

import (
	"fmt"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type Metrics struct {
	mu               sync.RWMutex
	startTime        time.Time
	requestCount     uint64
	successCount     uint64
	errorCount       uint64
	latencies        []time.Duration
	lastReportTime   time.Time
	lastRequestCount uint64
	heapAllocStart   uint64
	gcPauseStart     time.Duration
}

func NewMetrics() *Metrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &Metrics{
		startTime:      time.Now(),
		lastReportTime: time.Now(),
		heapAllocStart: m.HeapAlloc,
		gcPauseStart:   getDuration(m.PauseNs[(m.NumGC)%256]),
	}
}

func (m *Metrics) RecordLatency(latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	atomic.AddUint64(&m.requestCount, 1)
	atomic.AddUint64(&m.successCount, 1)
	m.latencies = append(m.latencies, latency)
}

func (m *Metrics) RecordError() {
	atomic.AddUint64(&m.requestCount, 1)
	atomic.AddUint64(&m.errorCount, 1)
}

func (m *Metrics) Report() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// elapsed := time.Since(m.startTime)
	currentCount := atomic.LoadUint64(&m.requestCount)
	// successCount := atomic.LoadUint64(&m.successCount)

	timeSinceLastReport := time.Since(m.lastReportTime)
	requestsSinceLastReport := currentCount - m.lastRequestCount
	rps := float64(requestsSinceLastReport) / timeSinceLastReport.Seconds()

	p50, p95, p99 := m.calculatePercentiles()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	goroutines := runtime.NumGoroutine()

	fmt.Printf("[%s] RPS: %.0f | Requests: %d | p50: %v | p95: %v | p99: %v | Goroutines: %d | Heap: %dMB | GC: %d\n",
		time.Now().Format("15:04:05"),
		rps,
		currentCount,
		p50,
		p95,
		p99,
		goroutines,
		memStats.HeapAlloc/1024/1024,
		memStats.NumGC,
	)
}

func (m *Metrics) calculatePercentiles() (p50, p95, p99 time.Duration) {
	if len(m.latencies) == 0 {
		return 0, 0, 0
	}

	sorted := make([]time.Duration, len(m.latencies))
	copy(sorted, m.latencies)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	length := len(sorted)
	p50 = sorted[length*50/100]
	p95 = sorted[length*95/100]
	p99 = sorted[length*99/100]

	return
}

func (m *Metrics) Final() {
	elapsed := time.Since(m.startTime)
	totalRequests := atomic.LoadUint64(&m.requestCount)
	successCount := atomic.LoadUint64(&m.successCount)
	errorCount := atomic.LoadUint64(&m.errorCount)
	avgRps := float64(totalRequests) / elapsed.Seconds()

	p50, p95, p99 := m.calculatePercentiles()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	fmt.Println("\n=== Final Summary ===")
	fmt.Printf("Duration: %v\n", elapsed)
	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Success: %d | Errors: %d\n", successCount, errorCount)
	fmt.Printf("Average RPS: %.0f\n", avgRps)
	fmt.Printf("Latency - p50: %v | p95: %v | p99: %v\n", p50, p95, p99)
	fmt.Printf("Heap: %dMB | GC Runs: %d\n", memStats.HeapAlloc/1024/1024, memStats.NumGC)
}

func getDuration(ns uint64) time.Duration {
	return time.Duration(ns) * time.Nanosecond
}
