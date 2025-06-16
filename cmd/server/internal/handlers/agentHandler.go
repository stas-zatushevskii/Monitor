package handlers

import (
	"fmt"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/constants"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/database"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/parser"
	"io"
	"net/http"
	"strconv"
)

func UpdateAgentHandler(storage *database.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data, err := parser.ParseJsonMetrics(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		nameMetric := data.ID
		typeMetric := data.MType

		switch typeMetric {
		case constants.Gauge:
			dataM := data.Value
			if err != nil {
				http.Error(w, constants.ErrParseFloat, http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			storage.SetGauge(nameMetric, *dataM)
			return
		case constants.Counter:
			dataM := data.Delta
			if err != nil {
				http.Error(w, constants.ErrParseInt, http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			storage.SetCounter(nameMetric, *dataM)
			return
		default:
			http.Error(w, constants.ErrUnsupportedType, http.StatusBadRequest)
			return
		}
	}
}

func ValueAgentHandler(storage *database.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := parser.ParseJsonMetrics(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		nameMetric := data.ID
		typeMetric := data.MType
		switch typeMetric {
		case constants.Gauge:
			value, ok := storage.GetGauge(nameMetric)
			if !ok {
				http.Error(w, constants.ErrGaugeNotFound, http.StatusNotFound)
				return
			}
			response := strconv.FormatFloat(value, 'f', -1, 64)
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, response)
			return
		case constants.Counter:
			value, ok := storage.GetCounter(nameMetric)
			if !ok {
				http.Error(w, constants.ErrCounterNotFound, http.StatusNotFound)
				return
			}
			response := strconv.FormatInt(value, 10)
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, response)
			return
		default:
			http.Error(w, constants.ErrUnsupportedType, http.StatusBadRequest)
			return
		}
	}
}

func GetAllAgentHandlers(storage *database.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "============- Gauge values -============\n")
		for key, val := range storage.Gauge {
			fmt.Fprintf(w, "	%s: %v\n", key, val)
		}
		fmt.Fprintf(w, "============- Counter values -============\n")
		for key, val := range storage.Counter {
			fmt.Fprintf(w, "	%s: %v\n", key, val)
		}
	}
}
