package metrics

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/workerpool"
)

// startSysMetricsCollector: In n period of time collect system metrics (library - gopsutil), and sends it in workerPool
func startSysMetricsCollector(ctx context.Context, wp *workerpool.WorkerPool, opt MonitorOptions, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(opt.PollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if vm, err := mem.VirtualMemory(); err == nil {
					submitGauge(wp, opt.URL, opt.HashKey, "TotalMemory", float64(vm.Total))
					submitGauge(wp, opt.URL, opt.HashKey, "FreeMemory", float64(vm.Free))
				} else {
					fmt.Println("gopsutil mem error:", err)
				}

				if percents, err := cpu.Percent(0, true); err == nil {
					for i := range percents {
						idx := i + 1
						p := percents[i]
						name := fmt.Sprintf("CPUutilization1_%d", idx)
						submitGauge(wp, opt.URL, opt.HashKey, name, p)
					}
				} else {
					fmt.Println("gopsutil cpu error:", err)
				}
			}
		}
	}()
}
