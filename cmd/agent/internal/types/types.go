package types

import "runtime"

type Gauge struct {
	Name string
	Data float64
}

type Counter struct {
	Name string
	Data int64
}

var GaugeMetrics = map[string]func(runtime.MemStats) float64{
	"Alloc":         func(m runtime.MemStats) float64 { return float64(m.Alloc) },
	"BuckHashSys":   func(m runtime.MemStats) float64 { return float64(m.BuckHashSys) },
	"Frees":         func(m runtime.MemStats) float64 { return float64(m.Frees) },
	"GCCPUFraction": func(m runtime.MemStats) float64 { return m.GCCPUFraction },
	"GCSys":         func(m runtime.MemStats) float64 { return float64(m.GCSys) },
	"HeapAlloc":     func(m runtime.MemStats) float64 { return float64(m.HeapAlloc) },
	"HeapIdle":      func(m runtime.MemStats) float64 { return float64(m.HeapIdle) },
	"HeapInuse":     func(m runtime.MemStats) float64 { return float64(m.HeapInuse) },
	"HeapObjects":   func(m runtime.MemStats) float64 { return float64(m.HeapObjects) },
	"HeapReleased":  func(m runtime.MemStats) float64 { return float64(m.HeapReleased) },
	"HeapSys":       func(m runtime.MemStats) float64 { return float64(m.HeapSys) },
	"LastGC":        func(m runtime.MemStats) float64 { return float64(m.LastGC) },
	"Lookups":       func(m runtime.MemStats) float64 { return float64(m.Lookups) },
	"MCacheInuse":   func(m runtime.MemStats) float64 { return float64(m.MCacheInuse) },
	"MCacheSys":     func(m runtime.MemStats) float64 { return float64(m.MCacheSys) },
	"MSpanInuse":    func(m runtime.MemStats) float64 { return float64(m.MSpanInuse) },
	"MSpanSys":      func(m runtime.MemStats) float64 { return float64(m.MSpanSys) },
	"Mallocs":       func(m runtime.MemStats) float64 { return float64(m.Mallocs) },
	"NextGC":        func(m runtime.MemStats) float64 { return float64(m.NextGC) },
	"NumForcedGC":   func(m runtime.MemStats) float64 { return float64(m.NumForcedGC) },
	"NumGC":         func(m runtime.MemStats) float64 { return float64(m.NumGC) },
	"OtherSys":      func(m runtime.MemStats) float64 { return float64(m.OtherSys) },
	"PauseTotalNs":  func(m runtime.MemStats) float64 { return float64(m.PauseTotalNs) },
	"StackInuse":    func(m runtime.MemStats) float64 { return float64(m.StackInuse) },
	"StackSys":      func(m runtime.MemStats) float64 { return float64(m.StackSys) },
	"Sys":           func(m runtime.MemStats) float64 { return float64(m.Sys) },
	"TotalAlloc":    func(m runtime.MemStats) float64 { return float64(m.TotalAlloc) },
}

var CounterMetrics = map[string]int64{
	"PollCount": 1,
}
