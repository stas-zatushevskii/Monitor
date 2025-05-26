package database

import (
	"fmt"
	"strings"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (ms *MemStorage) SetGauge(name string, data float64) {
	ms.gauge[name] = data
}

func (ms *MemStorage) SetCounter(name string, data int64) {
	ms.counter[name] += data
}

func (ms *MemStorage) GetGauge(name string) float64 {
	return ms.gauge[name]
}

func (ms *MemStorage) GetCounter(name string) int64 {
	return ms.counter[name]
}

func ParseData(url string) (nameM, dataM, typeM string, err error) {
	URLData := strings.Split(url, "/")

	if len(URLData) < 5 {
		return "", "", "", fmt.Errorf("invalid URL path: %v", url)
	}

	nameMetric := URLData[3]
	if nameMetric == "" {
		return "", "", "", fmt.Errorf("missing metric name in URL")
	}

	dataMetric := URLData[4]
	typeMetric := URLData[2]

	return nameMetric, dataMetric, typeMetric, nil
}
