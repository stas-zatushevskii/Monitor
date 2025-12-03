package metrics

import (
	"context"
	"fmt"
	"sync"

	"github.com/stas-zatushevskii/Monitor/cmd/agent/config"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/workerpool"
	"github.com/stas-zatushevskii/Monitor/cmd/proto"
)

func Monitor(ctx context.Context, url string, poolInterval, reportInterval int, cfg *config.Config, c proto.MetricsClient) {
	opt := NewMonitorOptions(url, poolInterval, reportInterval, cfg.RateLimit, 5, cfg.HashKey)

	wp := workerpool.NewWorkerPool(ctx, opt.RateLimit)
	defer wp.Close()

	var wg sync.WaitGroup
	store := &MemStatsStore{}

	startRuntimeBatchCollector(ctx, store, wp, opt, &wg, c)
	startPeriodicSender(ctx, store, wp, opt, &wg, c)
	startSysMetricsCollector(ctx, wp, opt, &wg, c)

	<-ctx.Done()
	fmt.Println("Agent stopping...")

	wg.Wait()
	wp.Close()
	fmt.Println("Agent stopped")
}
