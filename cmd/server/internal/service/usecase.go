package service

import "github.com/stas-zatushevskii/Monitor/cmd/server/internal/storage"

type MetricsService struct {
	storage storage.Storage
}

func NewMetricsService(storage storage.Storage) *MetricsService {
	return &MetricsService{storage: storage}
}
