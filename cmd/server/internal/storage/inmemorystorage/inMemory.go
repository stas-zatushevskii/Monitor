package inmemorystorage

import (
	"context"
	"fmt"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/models"
	"strings"
	"sync"
)

type InMemoryStorage struct {
	mu      sync.RWMutex
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		mu:      sync.RWMutex{},
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

func KeyToLower(key string) string {
	return strings.ToLower(key)
}

func (ms *InMemoryStorage) Snapshot() InMemoryStorage {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	gaugeCopy := make(map[string]float64, len(ms.Gauge))
	for k, v := range ms.Gauge {
		gaugeCopy[k] = v
	}

	counterCopy := make(map[string]int64, len(ms.Counter))
	for k, v := range ms.Counter {
		counterCopy[k] = v
	}

	return InMemoryStorage{
		Gauge:   gaugeCopy,
		Counter: counterCopy,
	}
}

func (ms *InMemoryStorage) SetGauge(name string, data float64) error {
	ms.Gauge[KeyToLower(name)] = data
	return nil
}

func (ms *InMemoryStorage) SetCounter(name string, data int64) error {
	ms.Counter[KeyToLower(name)] = data
	return nil
}

func (ms *InMemoryStorage) GetGauge(name string) (float64, error) {
	value, ok := ms.Gauge[KeyToLower(name)]
	if !ok {
		return 0, fmt.Errorf("gauge not found: %s", KeyToLower(name))
	}
	return value, nil
}

func (ms *InMemoryStorage) GetCounter(name string) (int64, error) {
	value, ok := ms.Counter[KeyToLower(name)]
	if !ok {
		return 0, fmt.Errorf("counter not found: %s", KeyToLower(name))
	}
	return value, nil
}

func (ms *InMemoryStorage) GetAllGauge() (map[string]float64, error) {
	return ms.Gauge, nil
}

func (ms *InMemoryStorage) GetAllCounter() (map[string]int64, error) {
	return ms.Counter, nil
}

// костыль для того чтобы имплементить интерфейс storage

func (ms *InMemoryStorage) Ping() error                         { return nil }
func (ms *InMemoryStorage) Close() error                        { return nil }
func (ms *InMemoryStorage) Bootstrap(ctx context.Context) error { return nil }
func (ms *InMemoryStorage) SetMultipleGauge(ctx context.Context, metrics []models.Metrics) error {
	return nil
}
func (ms *InMemoryStorage) SetMultipleCounter(ctx context.Context, metrics []models.Metrics) error {
	return nil
}
