package service

import (
	pb "github.com/stas-zatushevskii/Monitor/cmd/proto"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/storage"
)

type MetricsService struct {
	pb.UnimplementedMetricsServer  // auto generated gRPC methods
	storage storage.Storage
	hashKey string
}

func NewMetricsService(storage storage.Storage, hashKey string) *MetricsService {
	return &MetricsService{storage: storage, hashKey: hashKey}
}
