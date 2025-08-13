package metrics

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/stas-zatushevskii/Monitor/cmd/agent/config"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/sender"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/types"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/utils"

	"runtime"
)

type task struct {
	desc string
	fn   func() error
}

type workerPool struct {
	wg    sync.WaitGroup
	tasks chan task
}

func newWorkerPool(ctx context.Context, concurrency int) *workerPool {
	wp := &workerPool{
		tasks: make(chan task, 1024),
	}

	if concurrency <= 0 {
		concurrency = 1
	}

	wp.wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(id int) {
			defer wp.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case t, ok := <-wp.tasks:
					if !ok {
						return
					}
					if err := t.fn(); err != nil {
						fmt.Println("[worker", id, "] send error:", t.desc, ":", err)
					}
				}
			}
		}(i + 1)
	}
	return wp
}

func (wp *workerPool) submit(t task) {
	wp.tasks <- t
}

func (wp *workerPool) close() {
	close(wp.tasks)
	wp.wg.Wait()
}

func Monitor(ctx context.Context, url string, pollInterval, reportInterval int, cfg *config.Config) {

	wp := newWorkerPool(ctx, cfg.RateLimit)
	defer wp.close()

	tickerPoll := time.NewTicker(time.Duration(pollInterval) * time.Second)
	tickerSend := time.NewTicker(time.Duration(reportInterval) * time.Second)
	defer tickerPoll.Stop()
	defer tickerSend.Stop()

	var (
		m         runtime.MemStats
		tickCount int
	)
	const batchSize = 5
	var gaugeBuf []types.Gauge
	var counterBuf []types.Counter

	runtimeCollectorDone := make(chan struct{})
	go func() {
		defer close(runtimeCollectorDone)
		for {
			select {
			case <-ctx.Done():
				return
			case <-tickerPoll.C:
				runtime.ReadMemStats(&m)

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
				if tickCount >= batchSize {
					if len(gaugeBuf) > 0 {
						batch := append([]types.Gauge(nil), gaugeBuf...)
						wp.submit(task{
							desc: "gauge batch",
							fn: func() error {
								return sender.SendBatchData(batch, url)
							},
						})
					}
					if len(counterBuf) > 0 {
						batch := append([]types.Counter(nil), counterBuf...)
						wp.submit(task{
							desc: "counter batch",
							fn: func() error {
								return sender.SendBatchData(batch, url)
							},
						})
					}
					gaugeBuf = nil
					counterBuf = nil
					tickCount = 0
				}
			}
		}
	}()

	periodicSenderDone := make(chan struct{})
	go func() {
		defer close(periodicSenderDone)
		for {
			select {
			case <-ctx.Done():
				return
			case <-tickerSend.C:
				for name, fn := range types.GaugeMetrics {
					g := types.Gauge{Data: fn(m), Name: name}
					nameCopy := name
					wp.submit(task{
						desc: "gauge " + nameCopy,
						fn: func() error {
							return utils.RetryRequest(sender.SendData, g, url, cfg.HashKey)
						},
					})
				}
				for name, value := range types.CounterMetrics() {
					c := types.Counter{Data: value, Name: name}
					nameCopy := name
					wp.submit(task{
						desc: "counter " + nameCopy,
						fn: func() error {
							return utils.RetryRequest(sender.SendData, c, url, cfg.HashKey)
						},
					})
				}
			}
		}
	}()

	gopsutilDone := make(chan struct{})
	go func() {
		defer close(gopsutilDone)

		tickerSys := time.NewTicker(time.Duration(pollInterval) * time.Second)
		defer tickerSys.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-tickerSys.C:
				if vm, err := mem.VirtualMemory(); err == nil {
					total := types.Gauge{Name: "TotalMemory", Data: float64(vm.Total)}
					wp.submit(task{
						desc: "gauge TotalMemory",
						fn: func() error {
							return utils.RetryRequest(sender.SendData, total, url, cfg.HashKey)
						},
					})

					free := types.Gauge{Name: "FreeMemory", Data: float64(vm.Free)}
					wp.submit(task{
						desc: "gauge FreeMemory",
						fn: func() error {
							return utils.RetryRequest(sender.SendData, free, url, cfg.HashKey)
						},
					})
				} else {
					fmt.Println("gopsutil mem error:", err)
				}

				if percents, err := cpu.Percent(0, true); err == nil {
					for i, p := range percents {
						name := fmt.Sprintf("CPUutilization1_%d", i+1) // 1..N
						val := types.Gauge{Name: name, Data: p}
						wp.submit(task{
							desc: "gauge " + name,
							fn: func() error {
								return utils.RetryRequest(sender.SendData, val, url, cfg.HashKey)
							},
						})
					}
				} else {
					fmt.Println("gopsutil cpu error:", err)
				}
			}
		}
	}()

	<-ctx.Done()
	fmt.Println("Agent stopping...")

	<-runtimeCollectorDone
	<-periodicSenderDone
	<-gopsutilDone

	wp.close()
	fmt.Println("Agent stopped")
}
