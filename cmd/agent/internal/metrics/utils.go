package metrics

import (
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/sender"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/types"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/utils"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/workerpool"
	"runtime"
	"sync"
	"time"
)

type MemStatsStore struct {
	mu sync.RWMutex
	m  runtime.MemStats
}

func (s *MemStatsStore) Set(m runtime.MemStats) {
	s.mu.Lock()
	s.m = m
	s.mu.Unlock()
}

func (s *MemStatsStore) Get() runtime.MemStats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.m
}

type MonitorOptions struct {
	URL            string
	PollInterval   time.Duration // seconds
	ReportInterval time.Duration // seconds
	BatchSize      int
	HashKey        string
	RateLimit      int
}

func NewMonitorOptions(url string, pollIntervalSec, reportIntervalSec, rateLimit, batchSize int, hashKey string) MonitorOptions {
	if batchSize <= 0 {
		batchSize = 5
	}
	return MonitorOptions{
		URL:            url,
		PollInterval:   time.Duration(pollIntervalSec) * time.Second,
		ReportInterval: time.Duration(reportIntervalSec) * time.Second,
		BatchSize:      batchSize,
		HashKey:        hashKey,
		RateLimit:      rateLimit,
	}
}

// Functions for correctly send data in workerPool for each type of it

func submitGaugeBatch(wp *workerpool.WorkerPool, url string, batch []types.Gauge) {
	if len(batch) == 0 {
		return
	}
	payload := append([]types.Gauge(nil), batch...)
	wp.Submit(workerpool.Task{
		Desc: "gauge batch",
		Fn: func() error {
			return sender.SendBatchData(payload, url)
		},
	})
}

func submitCounterBatch(wp *workerpool.WorkerPool, url string, batch []types.Counter) {
	if len(batch) == 0 {
		return
	}
	payload := append([]types.Counter(nil), batch...)
	wp.Submit(workerpool.Task{
		Desc: "counter batch",
		Fn: func() error {
			return sender.SendBatchData(payload, url)
		},
	})
}

func submitGauge(wp *workerpool.WorkerPool, url, hashKey, name string, value float64) {
	g := types.Gauge{Name: name, Data: value}
	wp.Submit(workerpool.Task{
		Desc: "gauge " + name,
		Fn: func() error {
			return utils.RetryRequest(sender.SendData, g, url, hashKey)
		},
	})
}

func submitCounter(wp *workerpool.WorkerPool, url, hashKey, name string, value int64) {
	c := types.Counter{Name: name, Data: value}
	wp.Submit(workerpool.Task{
		Desc: "counter " + name,
		Fn: func() error {
			return utils.RetryRequest(sender.SendData, c, url, hashKey)
		},
	})
}
