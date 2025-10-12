package metrics

import (
	"context"
	"fmt"
	"sync"

	"github.com/stas-zatushevskii/Monitor/cmd/agent/config"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/workerpool"
)

func Monitor(ctx context.Context, url string, pollInterval, reportInterval int, cfg *config.Config) {
	opt := NewMonitorOptions(url, pollInterval, reportInterval, cfg.RateLimit, 5, cfg.HashKey)

	wp := workerpool.NewWorkerPool(ctx, opt.RateLimit)
	defer wp.Close()

	var wg sync.WaitGroup
	store := &MemStatsStore{}

	startRuntimeBatchCollector(ctx, store, wp, opt, &wg)
	startPeriodicSender(ctx, store, wp, opt, &wg)
	startSysMetricsCollector(ctx, wp, opt, &wg)

	<-ctx.Done()
	fmt.Println("Agent stopping...")

	wg.Wait()
	wp.Close()
	fmt.Println("Agent stopped")
}
