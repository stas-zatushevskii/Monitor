package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/stas-zatushevskii/Monitor/cmd/server/config"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/database"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/logger"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/router"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	config.ParseFlags()

	storage := database.NewMemStorage()
	if config.Restore {
		if err := database.AutoLoadData(config.FileStoragePath, storage); err != nil {
			log.Printf("Ошибка восстановления данных: %v", err)
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go database.AutoSaveData(ctx, storage, config.StoreInterval, config.FileStoragePath)

	r := router.New(storage)

	if err := run(r); err != nil {
		log.Fatal(err)
	}
}

func run(r *chi.Mux) error {
	if err := logger.Initialize(config.FlagLogLevel); err != nil {
		return err
	}
	logger.Log.Info("Running server", zap.String("address", config.Address))
	return http.ListenAndServe(config.Address, logger.WithLogging(r))
}
