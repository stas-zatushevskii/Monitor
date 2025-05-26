package handlers

import (
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/database"
	"net/http"
	"strconv"
)

func AgentHandler(storage *database.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(w, "Content-Type not allowed", http.StatusUnsupportedMediaType)
			return
		}
		if err := r.ParseForm(); err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// --- Функция получения данных из url
		nameMetric, dataMetric, typeMetric, err := database.ParseData(r.URL.String())
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		switch typeMetric {
		case "gauge":
			parsedData, err := strconv.ParseFloat(dataMetric, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			storage.SetGauge(nameMetric, parsedData)
			w.WriteHeader(http.StatusOK)
			return
		case "counter":
			parsedData, err := strconv.ParseInt(dataMetric, 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			storage.SetCounter(nameMetric, parsedData)
			w.WriteHeader(http.StatusOK)
			return
		default:
			http.Error(w, "Unsupported type", http.StatusBadRequest)
			return
		}
	}
}
