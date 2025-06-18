package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/stas-zatushevskii/Monitor/cmd/server/config"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/database"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/gzip"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/logger"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/router"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	storage := database.NewMemStorage()
	r := router.New(storage)
	r.Use(gzip.GzipMiddleware)
	config.ParseFlags()
	log.Fatal(run(r))
}

func run(r *chi.Mux) error {
	if err := logger.Initialize(config.FlagLogLevel); err != nil {
		return err
	}
	logger.Log.Info("Running server", zap.String("address", config.Address))
	return http.ListenAndServe(config.Address, logger.WithLogging(r))
}
