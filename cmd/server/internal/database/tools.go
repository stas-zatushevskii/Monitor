package database

import (
	"encoding/json"
	"errors"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/constants"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/parser"
	"strconv"
)

func SetJsonData(data parser.Metrics, storage Storage) error {
	switch data.MType {
	case constants.Gauge:
		dataM := data.Value
		storage.SetGauge(data.ID, *dataM)
		return nil
	case constants.Counter:
		dataM := data.Delta
		storage.SetCounter(data.ID, *dataM)
		return nil
	default:
		return errors.New(constants.ErrUnsupportedType)
	}
}

func SetURLData(nameMetric, dataMetric, typeMetric string, storage Storage) error {
	switch typeMetric {
	case constants.Gauge:
		parsedData, err := strconv.ParseFloat(dataMetric, 64)
		if err != nil {
			return errors.New(constants.ErrParseFloat)
		}
		storage.SetGauge(nameMetric, parsedData)
		return nil
	case constants.Counter:
		parsedData, err := strconv.ParseInt(dataMetric, 10, 64)
		if err != nil {
			return errors.New(constants.ErrParseFloat)
		}
		storage.SetCounter(nameMetric, parsedData)
		return nil
	default:
		return errors.New(constants.ErrUnsupportedType)
	}
}

func GetData(nameMetric, typeMetric string, storage Storage) (string, error) {
	switch typeMetric {
	case constants.Gauge:
		value, ok := storage.GetGauge(nameMetric)
		if !ok {
			return "", errors.New(constants.ErrGaugeNotFound)
		}
		response := strconv.FormatFloat(value, 'f', -1, 64)
		return response, nil
	case constants.Counter:
		value, ok := storage.GetCounter(nameMetric)
		if !ok {
			return "", errors.New(constants.ErrCounterNotFound)
		}
		response := strconv.FormatInt(value, 10)
		return response, nil
	default:
		return "", errors.New(constants.ErrUnsupportedType)
	}
}

func CreateJson(data Metrics) ([]byte, error) {
	result, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return result, nil
}
