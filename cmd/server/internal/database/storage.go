package database

import (
	"strings"
	"sync"
)

func KeyToLower(key string) string {
	return strings.ToLower(key)
}

type CounterStorage interface {
	SetCounter(name string, data int64)
	GetCounter(name string) (int64, bool)
}

type GaugeStorage interface {
	SetGauge(name string, data float64)
	GetGauge(name string) (float64, bool)
}

type Storage interface {
	CounterStorage
	GaugeStorage
}

type MemStorage struct {
	mu      sync.RWMutex
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

func (ms *MemStorage) SetGauge(name string, data float64) {
	ms.Gauge[KeyToLower(name)] = data
}

func (ms *MemStorage) SetCounter(name string, data int64) {
	ms.Counter[KeyToLower(name)] += data
}

func (ms *MemStorage) GetGauge(name string) (float64, bool) {
	value, ok := ms.Gauge[KeyToLower(name)]
	return value, ok
}

func (ms *MemStorage) GetCounter(name string) (int64, bool) {
	value, ok := ms.Counter[KeyToLower(name)]
	return value, ok
}

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
