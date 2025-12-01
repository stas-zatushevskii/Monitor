package metrics

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/sender"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/types"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/workerpool"
	"github.com/stas-zatushevskii/Monitor/cmd/proto"
)

// startSysMetricsCollector: In n period of time collect system metrics (library - gopsutil), and sends it in workerPool
func startSysMetricsCollector(ctx context.Context, wp *workerpool.WorkerPool, opt MonitorOptions, wg *sync.WaitGroup, c proto.MetricsClient) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(opt.PoolInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if vm, err := mem.VirtualMemory(); err == nil {
					submitData(wp, opt.URL, opt.HashKey, []types.Gauge{{Name: "TotalMemory", Data: float64(vm.Total)}})
					submitData(wp, opt.URL, opt.HashKey, []types.Gauge{{Name: "FreeMemory", Data: float64(vm.Free)}})
					if err := sender.SendgRPCData(context.Background(), []types.Gauge{{Name: "TotalMemory", Data: float64(vm.Total)}}, c); err != nil {
						return
					}
					if err := sender.SendgRPCData(context.Background(), []types.Gauge{{Name: "FreeMemory", Data: float64(vm.Free)}}, c); err != nil {
						return
					}
				} else {
					fmt.Println("gopsutil mem error:", err)
				}

				if percents, err := cpu.Percent(0, true); err == nil {
					for i := range percents {
						idx := i + 1
						p := percents[i]
						name := fmt.Sprintf("CPUutilization1_%d", idx)
						submitData(wp, opt.URL, opt.HashKey, []types.Gauge{{Name: name, Data: p}})
						if err := sender.SendgRPCData(context.Background(), []types.Gauge{{Name: name, Data: p}}, c); err != nil {
							return
						}
					}
				} else {
					fmt.Println("gopsutil cpu error:", err)
				}
			}
		}
	}()
}
