package main

import (
	"fmt"
	_ "net/http/pprof"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/api"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/audit"

	"github.com/stas-zatushevskii/Monitor/cmd/server/config"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/logger"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/service"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/storage"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/storage/inmemorystorage"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/storage/sqlstorage"

	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func printVersion() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	fmt.Printf("Build version: %s", buildVersion)
	fmt.Printf("Build date: %s", buildDate)
	fmt.Printf("Build commit: %s", buildCommit)
}

func main() {
	printVersion()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	// database
	var storage storage.Storage
	if cfg.DSN != "" {
		db, err := sql.Open("pgx", cfg.DSN)
		if err != nil {
			log.Fatalf("failed to connect to db: %v", err)
		}
		storage = sqlstorage.NewPostgresStorage(db)
		if err := storage.Bootstrap(ctx); err != nil {
			log.Fatalf("failed to bootstrap DB: %v", err)
		}
	} else {
		mem := inmemorystorage.NewInMemoryStorage()
		if cfg.Restore {
			_ = inmemorystorage.AutoLoadData(cfg.FileStoragePath, mem)
		}
		go inmemorystorage.AutoSaveData(ctx, mem, cfg.StoreInterval, cfg.FileStoragePath)
		storage = mem
	}

	// Service (depends on cfg what db it's use)
	metricsService := service.NewMetricsService(storage, cfg.HashKey)

	// audit
	logProducer := audit.NewLogProducer()
	logConsumer := audit.NewLogConsumer(cfg)
	logProducer.Register(logConsumer)

	//transport
	r := api.New(metricsService, cfg, logProducer)

	if err := run(r, ctx, cfg); err != nil {
		log.Fatal(err)
	}
	<-ctx.Done()
}

func run(r *chi.Mux, ctx context.Context, cfg *config.Config) error {
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		return err
	}
	logger.Log.Info("Running server", zap.String("address", cfg.Address))
	srv := &http.Server{
		Addr:    cfg.Address,
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
