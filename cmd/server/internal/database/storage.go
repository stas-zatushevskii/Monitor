package database

import (
	"fmt"
	"strings"
)

func KeyToLower(key string) string {
	return strings.ToLower(key)
}

type MemStorage struct {
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
