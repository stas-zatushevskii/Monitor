package service

import (
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/constants"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/models"

	"errors"
	"strconv"
)

func (m *MetricsService) SetJSONData(data models.Metrics) error {
	switch data.MType {
	case constants.Gauge:
		dataM := data.Value
		m.storage.SetGauge(data.ID, *dataM)
		return nil
	case constants.Counter:
		dataM := data.Delta
		m.storage.SetCounter(data.ID, *dataM)
		return nil
	default:
		return errors.New(constants.ErrUnsupportedType)
	}
}

func (m *MetricsService) SetURLData(nameMetric, dataMetric, typeMetric string) error {
	switch typeMetric {
	case constants.Gauge:
		parsedData, err := strconv.ParseFloat(dataMetric, 64)
		if err != nil {
			return errors.New(constants.ErrParseFloat)
		}
		m.storage.SetGauge(nameMetric, parsedData)
		return nil
	case constants.Counter:
		parsedData, err := strconv.ParseInt(dataMetric, 10, 64)
		if err != nil {
			return errors.New(constants.ErrParseFloat)
		}
		m.storage.SetCounter(nameMetric, parsedData)
		return nil
	default:
		return errors.New(constants.ErrUnsupportedType)
	}
}
