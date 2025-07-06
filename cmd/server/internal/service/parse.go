package service

import (
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/constants"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/models"

	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (m *MetricsService) ParseToMetrics(nameMetric, dataMetric, typeMetric string) (models.Metrics, error) {
	var data models.Metrics
	switch typeMetric {
	case constants.Gauge:
		parsedData, err := strconv.ParseFloat(dataMetric, 64)
		if err != nil {
			return data, err
		}
		data = models.Metrics{
			ID:    nameMetric,
			MType: typeMetric,
			Value: &parsedData,
		}
	case constants.Counter:
		parsedData, err := strconv.ParseInt(dataMetric, 10, 64)
		if err != nil {
			return data, err
		}
		data = models.Metrics{
			ID:    nameMetric,
			MType: typeMetric,
			Delta: &parsedData,
		}
	default:
		return data, errors.New(constants.ErrUnsupportedType)
	}
	return data, nil
}

func (m *MetricsService) ParseJSONData(r *http.Request) (models.Metrics, error) {
	var data models.Metrics
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return models.Metrics{}, fmt.Errorf(constants.ErrorParseJSON)
	}
	defer r.Body.Close()
	return data, nil
}

func (m *MetricsService) ParseJSONBatchData(r *http.Request) ([]models.Metrics, error) {
	var data []models.Metrics
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrorParseJSON)
	}
	defer r.Body.Close()
	return data, nil
}

func (m *MetricsService) ParseTypeMetrics(data []models.Metrics) (gauge, counter []models.Metrics, err error) {
	for _, metric := range data {
		switch metric.MType {
		case constants.Gauge:
			gauge = append(gauge, metric)
		case constants.Counter:
			counter = append(counter, metric)
		}
	}
	return gauge, counter, nil
}
