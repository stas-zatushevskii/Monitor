package metrics

import (
	"context"
	"sync"
	"time"

	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/sender"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/types"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/workerpool"
	"github.com/stas-zatushevskii/Monitor/cmd/proto"
)

// startPeriodicSender: start goroutine that in N period of time sends metrics in workerPool
func startPeriodicSender(ctx context.Context, store *MemStatsStore, wp *workerpool.WorkerPool, opt MonitorOptions, wg *sync.WaitGroup, c proto.MetricsClient) {
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
					submitData(wp, opt.URL, opt.HashKey, []types.Gauge{{Name: name, Data: val}})
					if err := sender.SendgRPCData(context.Background(), []types.Gauge{{Name: name, Data: val}}, c); err != nil {
						return
					}
				}
				for name, value := range types.CounterMetrics() {
					submitData(wp, opt.URL, opt.HashKey, []types.Counter{{Name: name, Data: value}})
					if err := sender.SendgRPCData(context.Background(), []types.Counter{{Name: name, Data: value}}, c); err != nil {
						return
					}
				}
			}
		}
	}()
}
