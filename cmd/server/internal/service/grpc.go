package service

import (
	"context"
	"net"

	pb "github.com/stas-zatushevskii/Monitor/cmd/proto"
	"github.com/stas-zatushevskii/Monitor/cmd/server/config"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/logger"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/models"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// 1 parse batch data
// 2 get ipAdr from metadata
// 3 set batch data with retry
// 4 set status to response

func (m *MetricsService) UpdateMetrics(ctx context.Context, in *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	var data []models.Metrics
	metricsData := in.GetMetrics()
	for _, metric := range metricsData {
		switch metric.GetType() {
		case pb.Metric_GAUGE:
			data = append(data, models.Metrics{
				ID:    metric.GetId(),
				MType: "gauge",
				Value: proto.Float64(metric.GetValue()),
				Delta: proto.Int64(metric.GetDelta()),
			})
		case pb.Metric_COUNTER:
			data = append(data, models.Metrics{
				ID:    metric.GetId(),
				MType: "counter",
				Value: proto.Float64(metric.GetValue()),
				Delta: proto.Int64(metric.GetDelta()),
			})
		case pb.Metric_Unspecified:
			return nil, status.Error(codes.InvalidArgument, "metric type not specified")
		}
	}

	err := utils.RetryWithContext(ctx, m.SetBatchData, data)
	if err != nil {
		return nil, err
	}

	response := &pb.UpdateMetricsResponse{}
	return response, nil
}

func UnaryInterceptor(cfg *config.Config) grpc.UnaryServerInterceptor {
	if cfg.TrustedSubnet == "" {
		return func(
			ctx context.Context,
			req interface{},
			info *grpc.UnaryServerInfo,
			handler grpc.UnaryHandler,
		) (interface{}, error) {
			return handler(ctx, req)
		}
	}

	_, ipnet, err := net.ParseCIDR(cfg.TrustedSubnet)
	if err != nil {
		return func(
			ctx context.Context,
			req interface{},
			info *grpc.UnaryServerInfo,
			handler grpc.UnaryHandler,
		) (interface{}, error) {
			logger.Log.Error(err.Error())
			return nil, status.Error(codes.PermissionDenied, "invalid trusted subnet")
		}
	}

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		var ipaddr string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get("X-Real-IP")
			if len(values) > 0 {
				ipaddr = values[0]
			}
		}

		if ipaddr == "" {
			return nil, status.Error(codes.PermissionDenied, "missing X-Real-IP")
		}

		ip := net.ParseIP(ipaddr)
		if ip == nil {
			return nil, status.Error(codes.PermissionDenied, "invalid ip")
		}

		if !ipnet.Contains(ip) {
			return nil, status.Error(codes.PermissionDenied, "access denied")
		}

		return handler(ctx, req)
	}
}
