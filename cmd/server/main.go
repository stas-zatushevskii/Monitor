package main

import (
	"errors"
	"fmt"
	"net"
	_ "net/http/pprof"
	"sync"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	pb "github.com/stas-zatushevskii/Monitor/cmd/proto"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/api/rest"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/audit"
	"google.golang.org/grpc"

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

	// logger
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		return
	}

	// Service (depends on cfg what db it's use)
	metricsService := service.NewMetricsService(storage, cfg.HashKey)

	// audit
	logProducer := audit.NewLogProducer()
	logConsumer := audit.NewLogConsumer(cfg)
	logProducer.Register(logConsumer)

	//transport
	r := rest.New(metricsService, cfg, logProducer)

	// run servers
	if err := run(ctx, cfg, r, metricsService); err != nil {
		logger.Log.Fatal("server stopped with error", zap.Error(err))
	}
}

// run starting REST and gRPC server
func run(ctx context.Context, cfg *config.Config, r *chi.Mux, metricsService *service.MetricsService) error {
	errCh := make(chan error, 2)

	restServer := runREST(r, cfg, errCh)
	grpcServer := runGRPC(cfg, metricsService, errCh)

	select {
	case err := <-errCh:
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownServers(shutdownCtx, restServer, grpcServer)

		return err

	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownServers(shutdownCtx, restServer, grpcServer)

		return nil
	}
}

// shutdownServers asynchronously shut down servers
func shutdownServers(ctx context.Context, restServer *http.Server, grpcServer *grpc.Server) {
	var wg sync.WaitGroup

	// REST
	if restServer != nil {
		wg.Add(1)

		go func() {
			defer wg.Done()

			if err := restServer.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Log.Warn("REST shutdown error", zap.Error(err))
			}
		}()
	}

	// gRPC
	if grpcServer != nil {
		wg.Add(1)

		go func() {
			defer wg.Done()

			done := make(chan struct{})
			// try to close server gracefully
			go func() {
				grpcServer.GracefulStop()
				close(done)
			}()

			select {
			case <-done:
			case <-ctx.Done():
				// in case, then timeout exited
				// manually closing server
				grpcServer.Stop()
			}
		}()
	}
	wg.Wait()
}

// runREST runs REST server
func runREST(r *chi.Mux, cfg *config.Config, errCh chan<- error) *http.Server {
	logger.Log.Info("Running REST server", zap.String("address", cfg.Address))

	s := &http.Server{
		Addr:    cfg.Address,
		Handler: logger.WithLogging(r),
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("REST serve: %w", err)
		}
	}()

	return s
}

// runGRPC runs gRPC server
func runGRPC(cfg *config.Config, metricService *service.MetricsService, errCh chan<- error) *grpc.Server {
	logger.Log.Info("Running gRPC server", zap.String("address", cfg.AddressGRPC))

	lis, err := net.Listen("tcp", cfg.AddressGRPC)
	if err != nil {
		errCh <- fmt.Errorf("gRPC listen: %w", err)
		return nil
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(service.UnaryInterceptor(cfg)),
	)

	pb.RegisterMetricsServer(s, metricService)

	go func() {
		if err := s.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			errCh <- fmt.Errorf("gRPC serve: %w", err)
		}
	}()

	return s
}
