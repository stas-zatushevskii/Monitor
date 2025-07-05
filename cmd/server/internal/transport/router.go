package transport

import (
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/gzip"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/service"

	"github.com/go-chi/chi/v5"
)

func New(metricService *service.MetricsService) *chi.Mux {
	router := chi.NewRouter()
	handler := NewHandler(metricService)

	router.Use(gzip.GzipMiddleware)

	router.Get("/", handler.GetAllAgentHandlers())
	router.Post("/update/", handler.UpdateJSONHandler())
	router.Post("/value/", handler.ValueJSONHandler())
	router.Post("/update/{type}/{name}/{data}", handler.UpdateURLHandler())
	router.Get("/value/{type}/{name}", handler.ValueURLHandler())
	router.Get("/ping", handler.Ping())
	return router
}
