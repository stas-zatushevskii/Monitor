package service

import (
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/storage"
)

type MetricsService struct {
	storage storage.Storage
	hashKey string
}

func NewMetricsService(storage storage.Storage, hashKey string) *MetricsService {
	return &MetricsService{storage: storage, hashKey: hashKey}
}
