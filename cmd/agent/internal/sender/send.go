package sender

import (
	"encoding/json"
	"errors"

	"github.com/go-resty/resty/v2"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/hash"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/types"
)

// CreateMetrics
// Convert data to struct types.Metrics
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
		SetBody(metrics). // JSON сериализуется автоматически
		Post(updateURL)

	if hashKey != "" {
		hashData := hash.HashData(data, []byte(hashKey))
		client.R().SetHeader("HashSHA256", hashData)
	}
	return err
}
