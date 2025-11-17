package metrics

import (
	"runtime"
	"sync"
	"time"

	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/sender"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/types"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/workerpool"
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
	PoolInterval   time.Duration // seconds
	ReportInterval time.Duration // seconds
	BatchSize      int
	HashKey        string
	RateLimit      int
}

func NewMonitorOptions(url string, poolIntervalSec, reportIntervalSec, rateLimit, batchSize int, hashKey string) MonitorOptions {
	if batchSize <= 0 {
		batchSize = 5
	}
	return MonitorOptions{
		URL:            url,
		PoolInterval:   time.Duration(poolIntervalSec) * time.Second,
		ReportInterval: time.Duration(reportIntervalSec) * time.Second,
		BatchSize:      batchSize,
		HashKey:        hashKey,
		RateLimit:      rateLimit,
	}
}

// Functions for correctly send data in workerPool for each type of it

func submitData[D types.MetricData](wp *workerpool.WorkerPool, url, hashKey string, data []D) {
	wp.Submit(workerpool.Task{
		Desc: "gauge batch",
		Fn: func() error {
			return sender.SendData(data, url, hashKey)
		},
	})
}
