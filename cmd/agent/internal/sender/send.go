package sender

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

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

func SendData[metricData types.MetricData](m metricData, url, hashKey string) error {
	updateURL := url + "/update/"
	parsedMetric, err := CreateMetrics(m)
	if err != nil {
		return err
	}
	data, err := json.Marshal(parsedMetric)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest(http.MethodPost, updateURL, bytes.NewReader(data))
	if hashKey != "" {
		hashData := hash.HashData(data, []byte(hashKey))
		req.Header.Set("HashSHA256", hashData)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: time.Second * 2,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status code: %d", resp.StatusCode)
	}
	return nil
}

func SendBatchData[metricData types.MetricData](m []metricData, url string) error {
	updateURL := url + "/updates/"
	var metrics []types.Metrics

	for _, metric := range m {
		parsed, err := CreateMetrics(metric)
		if err != nil {
			return err
		}
		metrics = append(metrics, parsed)
	}

	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(metrics). // JSON сериализуется автоматически
		Post(updateURL)

	return err
}
