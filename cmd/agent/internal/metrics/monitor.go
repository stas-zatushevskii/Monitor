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
	for {
		select {
		case <-tickerPool.C:
			runtime.ReadMemStats(&m)
		case <-tickerSend.C:
			for name, fn := range types.GaugeMetrics {
				if err := sender.SendData(types.Gauge{Data: fn(m), Name: name}, url); err != nil {
					fmt.Println("Error sending metric:", err)
				}
			}
			for name, value := range types.CounterMetrics {
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
