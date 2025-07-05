package service

import (
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/constants"

	"errors"
	"strconv"
)

func (m *MetricsService) GetDataByName(nameMetric, typeMetric string) (string, error) {
	switch typeMetric {
	case constants.Gauge:
		value, err := m.storage.GetGauge(nameMetric)
		if err != nil {
			return "", errors.New(constants.ErrGaugeNotFound)
		}
		response := strconv.FormatFloat(value, 'f', -1, 64)
		return response, nil
	case constants.Counter:
		value, err := m.storage.GetCounter(nameMetric)
		if err != nil {
			return "", errors.New(constants.ErrCounterNotFound)
		}
		response := strconv.FormatInt(value, 10)
		return response, nil
	default:
		return "", errors.New(constants.ErrUnsupportedType)
	}
}

func (m *MetricsService) GetAllGaugeMetrics() (map[string]float64, error) {
	data, err := m.storage.GetAllGauge()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (m *MetricsService) GetAllCounterMetrics() (map[string]int64, error) {
	data, err := m.storage.GetAllCounter()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (m *MetricsService) GetDataBaseStatus() error {
	err := m.storage.Ping()
	return err
}
