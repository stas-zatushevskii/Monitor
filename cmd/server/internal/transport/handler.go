package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/constants"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/service"
)

type Handler struct {
	metricService *service.MetricsService
}

func NewHandler(metricService *service.MetricsService) *Handler {
	return &Handler{metricService: metricService}
}

func (h *Handler) UpdateJSONHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data, err := h.metricService.ParseJSONData(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.metricService.SetJSONData(data)
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

func (h *Handler) UpdateURLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nameMetric := chi.URLParam(r, "name")
		dataMetric := chi.URLParam(r, "data")
		typeMetric := chi.URLParam(r, "type")

		err := h.metricService.SetURLData(nameMetric, dataMetric, typeMetric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) ValueJSONHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := h.metricService.ParseJSONData(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		nameMetric := data.ID
		typeMetric := data.MType
		dataMetric, err := h.metricService.GetDataByName(nameMetric, typeMetric)
		if err != nil {
			if err.Error() == constants.ErrCounterNotFound || err.Error() == constants.ErrGaugeNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		metrics, _ := h.metricService.ParseToMetrics(nameMetric, dataMetric, typeMetric)
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

func (h *Handler) ValueURLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nameMetric := chi.URLParam(r, "name")
		typeMetric := chi.URLParam(r, "type")

		response, err := h.metricService.GetDataByName(nameMetric, typeMetric)
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

func (h *Handler) GetAllAgentHandlers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		gauge, err := h.metricService.GetAllGaugeMetrics()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "============- Gauge values -============\n")
		for key, val := range gauge {
			fmt.Fprintf(w, "	%s: %v\n", key, val)
		}

		counter, err := h.metricService.GetAllCounterMetrics()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "============- Counter values -============\n")
		for key, val := range counter {
			fmt.Fprintf(w, "	%s: %v\n", key, val)
		}
	}
}

func (h *Handler) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h.metricService.GetDataBaseStatus(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err.Error())
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
	}
}
