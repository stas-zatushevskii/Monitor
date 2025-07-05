package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/stas-zatushevskii/Monitor/cmd/server/config"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/logger"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/service"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/storage"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/storage/inmemorystorage"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/storage/sqlstorage"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/transport"

	"context"
	"database/sql"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type application struct {
	MetricsService *service.MetricsService
}

func main() {
	config.ParseFlags()
	var storage storage.Storage
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if config.DSN != "" {
		db, err := sql.Open("pgx", config.DSN)
		if err != nil {
			log.Fatalf("failed to connect to db: %v", err)
		}
		storage = sqlstorage.NewPostgresStorage(db)
	} else {
		mem := inmemorystorage.NewInMemoryStorage()
		if config.Restore {
			_ = inmemorystorage.AutoLoadData(config.FileStoragePath, mem)
		}
		go inmemorystorage.AutoSaveData(ctx, mem, config.StoreInterval, config.FileStoragePath)
		storage = mem
	}

	// Service (depends on cfg what db it's use)
	metricsService := service.NewMetricsService(storage)

	r := transport.New(metricsService)

	if err := run(r, ctx); err != nil {
		log.Fatal(err)
	}
	<-ctx.Done()
}

func run(r *chi.Mux, ctx context.Context) error {
	if err := logger.Initialize(config.FlagLogLevel); err != nil {
		return err
	}
	logger.Log.Info("Running server", zap.String("address", config.Address))
	srv := &http.Server{
		Addr:    config.Address,
		Handler: logger.WithLogging(r),
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()
	select {
	case <-ctx.Done():
		logger.Log.Info("Shutting down server...")
		return srv.Shutdown(context.Background())
	case err := <-errCh:
		return err
	}
}
