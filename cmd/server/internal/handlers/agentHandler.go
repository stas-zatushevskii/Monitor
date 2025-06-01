package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/database"
	"io"
	"net/http"
	"strconv"
)

func UpdateAgentHandler(storage *database.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		nameMetric := chi.URLParam(r, "name")
		dataMetric := chi.URLParam(r, "data")
		typeMetric := chi.URLParam(r, "type")

		switch typeMetric {
		case Gauge:
			parsedData, err := strconv.ParseFloat(dataMetric, 64)
			if err != nil {
				http.Error(w, ErrParseFloat, http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			storage.SetGauge(nameMetric, parsedData)
			return
		case Counter:
			parsedData, err := strconv.ParseInt(dataMetric, 10, 64)
			if err != nil {
				http.Error(w, ErrParseInt, http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			storage.SetCounter(nameMetric, parsedData)
			return
		default:
			http.Error(w, ErrUnsupportedType, http.StatusBadRequest)
			return
		}
	}
}

func ValueAgentHandler(storage *database.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nameMetric := chi.URLParam(r, "name")
		typeMetric := chi.URLParam(r, "type")
		switch typeMetric {
		case Gauge:
			value, ok := storage.GetGauge(nameMetric)
			if !ok {
				http.Error(w, ErrGaugeNotFound, http.StatusNotFound)
				return
			}
			response := strconv.FormatFloat(value, 'f', -1, 64)
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, response)
			return
		case Counter:
			value, ok := storage.GetCounter(nameMetric)
			if !ok {
				http.Error(w, ErrCounterNotFound, http.StatusNotFound)
				return
			}
			response := strconv.FormatInt(value, 10)
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, response)
			return
		default:
			http.Error(w, ErrUnsupportedType, http.StatusBadRequest)
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
