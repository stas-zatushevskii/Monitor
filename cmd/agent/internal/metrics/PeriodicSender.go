package metrics

import (
	"context"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/types"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/workerpool"
	"sync"
	"time"
)

// startPeriodicSender: start goroutine that in N period of time sends metrics in workerPool
func startPeriodicSender(ctx context.Context, store *MemStatsStore, wp *workerpool.WorkerPool, opt MonitorOptions, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(opt.ReportInterval)
		defer ticker.Stop()

		for {
			select {
			// listen GracefulShutdown
			case <-ctx.Done():
				return
			case <-ticker.C:
				m := store.Get()
				for name, fn := range types.GaugeMetrics {
					val := fn(m)
					submitGauge(wp, opt.URL, opt.HashKey, name, val)
				}
				for name, value := range types.CounterMetrics() {
					submitCounter(wp, opt.URL, opt.HashKey, name, value)
				}
			}
		}
	}()
}
