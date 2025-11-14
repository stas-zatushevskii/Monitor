package metrics

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/types"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/workerpool"
)

// startRuntimeBatchCollector: In n period of time collect batch of metrics each type, and sends it in workerPool
func startRuntimeBatchCollector(ctx context.Context, store *MemStatsStore, wp *workerpool.WorkerPool, opt MonitorOptions, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(opt.PoolInterval)
		defer ticker.Stop()

		var tickCount int
		var gaugeBuf []types.Gauge
		var counterBuf []types.Counter

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				var m runtime.MemStats
				runtime.ReadMemStats(&m)
				store.Set(m)

				// Fill buffers
				for name, fn := range types.GaugeMetrics {
					gaugeBuf = append(gaugeBuf, types.Gauge{
						Name: name,
						Data: fn(m),
					})
				}
				for name, value := range types.CounterMetrics() {
					counterBuf = append(counterBuf, types.Counter{
						Name: name,
						Data: value,
					})
				}

				tickCount++
				if tickCount >= opt.BatchSize {
					submitGaugeBatch(wp, opt.URL, gaugeBuf)
					submitCounterBatch(wp, opt.URL, counterBuf)
					gaugeBuf = nil
					counterBuf = nil
					tickCount = 0
				}
			}
		}
	}()
}
