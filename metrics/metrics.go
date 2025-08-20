package metrics

import (
	"sync/atomic"
	"time"

	"metrics/logger"
)

var txCount int64

func IncrementTPS() {
	atomic.AddInt64(&txCount, 1)
}

// StartTPSLogger logs TPS every second
func StartTPSLogger() {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			count := atomic.SwapInt64(&txCount, 0)
			if count > 0 {
				logger.Metrics.Printf("TPS = %d", count)
			}
		}
	}()
}

// LogLatency logs latency for an RPC call
func LogLatency(method string, duration time.Duration) {
	logger.Metrics.Printf("Latency for %s = %v", method, duration)
}

// LogFinality logs finality time for a tx
func LogFinality(txID string, duration time.Duration) {
	logger.Metrics.Printf("Time to finality for %s = %v", txID, duration)
}

