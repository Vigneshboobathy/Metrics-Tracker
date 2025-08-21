package metrics

import (
	"sync/atomic"
	"time"

	"metrics/logger"
)

var (
	txCount   int64
	tpsTicker *time.Ticker
	stopTPS   chan struct{}
)

// StartTPSLogger logs TPS every second
func StartTPSLogger() {
	if tpsTicker != nil {
		return
	}
	tpsTicker = time.NewTicker(1 * time.Second)
	stopTPS = make(chan struct{})
	go func() {
		for {
			select {
			case <-tpsTicker.C:
				count := atomic.SwapInt64(&txCount, 0)
				logger.Metrics.Printf("TPS = %d", count)
			case <-stopTPS:
				tpsTicker.Stop()
				tpsTicker = nil
				return
			}
		}
	}()
}

// LogLatency logs latency for an RPC call
func LogLatency(method string, duration time.Duration) {
	logger.Metrics.Printf("Latency for %s = %v", method, duration)
}
