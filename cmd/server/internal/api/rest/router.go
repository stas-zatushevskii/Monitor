package rest

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stas-zatushevskii/Monitor/cmd/server/config"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/audit"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/gzip"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/service"
	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
)

func New(metricService *service.MetricsService, config *config.Config, audit *audit.LogProducer, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()
	handler := NewHandler(metricService, config, audit)

	router.Use(gzip.GzipMiddleware)
	router.Use(WhiteListMiddleware(config))

	router.Mount("/debug", middleware.Profiler())

	router.Get("/", handler.GetAllAgentHandlers())
	router.Post("/update/", handler.UpdateJSONHandler()) //
	router.Post("/updates/", handler.SetBatchDataJSON()) //
	router.Post("/value/", handler.ValueJSONHandler())
	router.Post("/update/{type}/{name}/{data}", handler.UpdateURLHandler())
	router.Get("/value/{type}/{name}", handler.ValueURLHandler())
	router.Get("/ping", handler.Ping())
	return router
}
