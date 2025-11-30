package sender

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/hash"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/types"
	pb "github.com/stas-zatushevskii/Monitor/cmd/proto"
	"google.golang.org/grpc/metadata"
)

func CreateMetrics[metricData types.Gauge | types.Counter](m metricData) (types.Metrics, error) {
	switch v := any(m).(type) {
	case types.Counter:
		return types.Metrics{
			ID:    v.Name,
			MType: "counter",
			Delta: &v.Data,
		}, nil
	case types.Gauge:
		return types.Metrics{
			ID:    v.Name,
			MType: "gauge",
			Value: &v.Data,
		}, nil
	default:
		return types.Metrics{}, errors.New("invalid metric type")
	}
}

func SendData[metricData types.MetricData](m []metricData, url, hashKey string) error {
	updateURL := url + "/updates/"
	var metrics []types.Metrics

	for _, metric := range m {
		parsed, err := CreateMetrics(metric)
		if err != nil {
			return err
		}
		metrics = append(metrics, parsed)
	}

	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	client := resty.New()
	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Real-IP", "127.0.0.1").
		SetBody(metrics).
		Post(updateURL)

	if hashKey != "" {
		hashData := hash.HashData(data, []byte(hashKey))
		client.R().SetHeader("HashSHA256", hashData)
	}
	return err
}

func SendgRPCData[M types.MetricData](ctx context.Context, m []M, c pb.MetricsClient) error {
	metrics := make([]types.Metrics, 0, len(m))
	for _, metric := range m {
		parsed, err := CreateMetrics(metric)
		if err != nil {
			return fmt.Errorf("create metrics: %w", err)
		}
		metrics = append(metrics, parsed)
	}

	data := make([]*pb.Metric, 0, len(metrics))
	for _, metric := range metrics {
		var mType pb.Metric_MType
		switch metric.MType {
		case "counter":
			mType = pb.Metric_COUNTER
		case "gauge":
			mType = pb.Metric_GAUGE
		default:
			return fmt.Errorf("invalid metric type: %s", metric.MType)
		}

		var delta int64
		if metric.Delta != nil {
			delta = *metric.Delta
		}
		var value float64
		if metric.Value != nil {
			value = *metric.Value
		}

		pm := (&pb.Metric_builder{
			Id:    metric.ID,
			Type:  mType,
			Delta: delta,
			Value: value,
		}).Build()

		data = append(data, pm)
	}

	md := metadata.Pairs("X-Real-IP", "127.0.0.1")
	ctx = metadata.NewOutgoingContext(ctx, md)

	req := (&pb.UpdateMetricsRequest_builder{
		Metrics: data,
	}).Build()

	_, err := c.UpdateMetrics(ctx, req)
	if err != nil {
		return fmt.Errorf("update metrics: %w", err)
	}

	return nil
}
