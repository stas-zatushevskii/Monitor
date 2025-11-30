package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/stas-zatushevskii/Monitor/cmd/agent/config"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/metrics"
	pb "github.com/stas-zatushevskii/Monitor/cmd/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	// сигнал который мониторит принудительную остановку программы и отправляет в контекст
	// stop() все равно остановит программу если вдруг сигнал не отправится
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// запускается в фоне и будет остановлена только после принудительной остановки программы
	url := "http://" + cfg.Address

	// grpc
	conn, err := grpc.NewClient(cfg.AddressGRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("failed to close: %v", err)
		}
	}()

	c := pb.NewMetricsClient(conn)

	go metrics.Monitor(ctx, url, cfg.PollInterval, cfg.ReportInterval, cfg, c)
	<-ctx.Done()
}
