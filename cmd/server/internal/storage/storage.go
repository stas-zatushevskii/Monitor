package storage

import (
	"context"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/models"
)

// in memory storage

type CounterStorage interface {
	SetCounter(name string, data int64) error
	GetCounter(name string) (int64, error)
	GetAllCounter() (map[string]int64, error)
	SetMultipleCounter(ctx context.Context, metrics []models.Metrics) error
}
type GaugeStorage interface {
	SetGauge(name string, data float64) error
	GetGauge(name string) (float64, error)
	GetAllGauge() (map[string]float64, error)
	SetMultipleGauge(ctx context.Context, metrics []models.Metrics) error
}

// abstract storage

type Storage interface {
	CounterStorage
	GaugeStorage
	Ping() error
	Close() error
	Bootstrap(ctx context.Context) error
}
