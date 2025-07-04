package main

import (
	"context"
	"database/sql"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
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

	ps := config.DSN
	db, err := sql.Open("pgx", ps)
	if err != nil {
		log.Printf("failed to open database: %v", err)
	}
	defer db.Close()

	storage := database.NewMemStorage()
	if config.Restore {
		if err := database.AutoLoadData(config.FileStoragePath, storage); err != nil {
			log.Printf("Ошибка восстановления данных: %v", err)
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go database.AutoSaveData(ctx, storage, config.StoreInterval, config.FileStoragePath)

	r := router.New(storage, db)

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
