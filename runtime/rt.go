// Package runtime provides extensions to Go's standard runtime package.
// It offers additional utility functions for runtime information, management,
// and diagnostics beyond what's provided in the standard library.
package runtime

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"runtime/trace"
	"time"
)

// MemStats extends runtime.MemStats with additional computed fields and helper methods.
type MemStats struct {
	runtime.MemStats
	UsagePercent float64 // Memory usage as percentage of total system memory
}

// GetMemStats returns enhanced memory statistics.
func GetMemStats() *MemStats {
	var stats MemStats
	runtime.ReadMemStats(&stats.MemStats)

	// Calculate additional metrics
	stats.UsagePercent = float64(stats.Alloc) / float64(stats.Sys) * 100

	return &stats
}

// ForceGC forces an immediate garbage collection with optional waiting.
func ForceGC(wait bool) {
	runtime.GC()
	if wait {
		// Run another GC and wait to ensure completion
		debug.FreeOSMemory()
	}
}

// GoroutineStats provides information about goroutines.
type GoroutineStats struct {
	Count       int
	Idle        int
	Running     int
	BlockedOn   map[string]int
	StackTraces []string
}

// GetGoroutineStats returns statistics about all goroutines.
func GetGoroutineStats(includeTraces bool) (*GoroutineStats, error) {
	stats := &GoroutineStats{
		Count:     runtime.NumGoroutine(),
		BlockedOn: make(map[string]int),
	}

	if includeTraces {
		buf := make([]byte, 1<<20)
		n := runtime.Stack(buf, true)
		stats.StackTraces = []string{string(buf[:n])}
	}

	return stats, nil
}

// CPUProfileStart starts CPU profiling and writes to the given file.
// Returns a stop function that should be called to end profiling.
func CPUProfileStart(filename string) (func(), error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("could not create CPU profile: %v", err)
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		f.Close()
		return nil, fmt.Errorf("could not start CPU profile: %v", err)
	}

	return func() {
		pprof.StopCPUProfile()
		f.Close()
	}, nil
}

// TraceStart starts execution tracing and writes to the given file.
// Returns a stop function that should be called to end tracing.
func TraceStart(filename string) (func(), error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("could not create trace file: %v", err)
	}

	if err := trace.Start(f); err != nil {
		f.Close()
		return nil, fmt.Errorf("could not start tracing: %v", err)
	}

	return func() {
		trace.Stop()
		f.Close()
	}, nil
}

// WithTimeout runs the given function with a specified timeout.
// Returns true if function completed, false if it timed out.
func WithTimeout(timeout time.Duration, fn func()) bool {
	done := make(chan struct{})

	go func() {
		fn()
		close(done)
	}()

	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}

// SetMaxThreads sets the maximum number of OS threads that can be executing
// Go code simultaneously and returns the previous setting.
func SetMaxThreads(n int) int {
	prev := runtime.GOMAXPROCS(0)
	if n > 0 {
		runtime.GOMAXPROCS(n)
	}
	return prev
}
