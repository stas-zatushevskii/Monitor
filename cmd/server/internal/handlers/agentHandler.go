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
		case "gauge":
			parsedData, err := strconv.ParseFloat(dataMetric, 64)
			if err != nil {
				http.Error(w, "Error while parsing float", http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			storage.SetGauge(nameMetric, parsedData)
			return
		case "counter":
			parsedData, err := strconv.ParseInt(dataMetric, 10, 64)
			if err != nil {
				http.Error(w, "Error while parsing int", http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			storage.SetCounter(nameMetric, parsedData)
			return
		default:
			http.Error(w, "Unsupported type", http.StatusBadRequest)
			return
		}
	}
}

func ValueAgentHandler(storage *database.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nameMetric := chi.URLParam(r, "name")
		valueMetric := chi.URLParam(r, "value")
		switch valueMetric {
		case "gauge":
			value, ok := storage.GetGauge(nameMetric)
			if !ok {
				http.Error(w, "gauge not found", http.StatusNotFound)
				return
			}
			response := nameMetric + ": " + strconv.FormatFloat(value, 'f', -1, 64)
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, response)
			return
		case "counter":
			value, ok := storage.GetCounter(nameMetric)
			if !ok {
				http.Error(w, "counter not found", http.StatusNotFound)
				return
			}
			response := nameMetric + ": " + strconv.FormatInt(value, 10)
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, response)
			return
		default:
			http.Error(w, "Unsupported type", http.StatusBadRequest)
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
