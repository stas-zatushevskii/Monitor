package metrics

import (
	"context"
	"fmt"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/sender"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/types"
	"runtime"
	"time"
)

func Monitor(ctx context.Context, url string, pollInterval, reportInterval int) {
	tickerPool := time.NewTicker(time.Duration(pollInterval) * time.Second)
	tickerSend := time.NewTicker(time.Duration(reportInterval) * time.Second)
	defer tickerPool.Stop()
	defer tickerSend.Stop()
	var m runtime.MemStats
	var tickCount int
	const batchSize = 5

	var gaugeBuf []types.Gauge
	var counterBuf []types.Counter

	for {
		select {
		case <-tickerPool.C:
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
				// Отправка батчем
				if err := sender.SendBatchData(gaugeBuf, url); err != nil {
					fmt.Println("Gauge batch send error:", err)
				}
				if err := sender.SendBatchData(counterBuf, url); err != nil {
					fmt.Println("Counter batch send error:", err)
				}
				// Очистка
				gaugeBuf = nil
				counterBuf = nil
				tickCount = 0
			}
		case <-tickerSend.C:
			for name, fn := range types.GaugeMetrics {
				if err := sender.SendData(types.Gauge{Data: fn(m), Name: name}, url); err != nil {
					fmt.Println("Error sending metric:", err)
				}
			}
			for name, value := range types.CounterMetrics() {
				if err := sender.SendData(types.Counter{Data: value, Name: name}, url); err != nil {
					fmt.Println("Error sending metric:", err)
				}
			}

		case <-ctx.Done():
			fmt.Println("Agent stopped")
			return
		}
	}
}
