package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/constants"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/database"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/parser"
	"io"
	"net/http"
)

func UpdateJSONHandler(storage *database.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data, err := parser.ParseJSONData(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = database.SetJSONData(data, storage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}

func UpdateURLHandler(storage *database.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nameMetric := chi.URLParam(r, "name")
		dataMetric := chi.URLParam(r, "data")
		typeMetric := chi.URLParam(r, "type")

		err := database.SetURLData(nameMetric, dataMetric, typeMetric, storage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
	}
}

func ValueJSONHandler(storage *database.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := parser.ParseJSONData(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		nameMetric := data.ID
		typeMetric := data.MType
		dataMetric, err := database.GetData(nameMetric, typeMetric, storage)
		if err != nil {
			if err.Error() == constants.ErrCounterNotFound || err.Error() == constants.ErrGaugeNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		metrics, _ := database.ParseToMetrics(nameMetric, dataMetric, typeMetric)
		result, err := json.Marshal(metrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}

func ValueURLHandler(storage *database.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nameMetric := chi.URLParam(r, "name")
		typeMetric := chi.URLParam(r, "type")

		response, err := database.GetData(nameMetric, typeMetric, storage)
		if err != nil {
			if err.Error() == constants.ErrCounterNotFound || err.Error() == constants.ErrGaugeNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, response)
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
