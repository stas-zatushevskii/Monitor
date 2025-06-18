package main

import (
	"context"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/config"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/metrics"
	"log"
	"os"
	"os/signal"
)

func main() {
	// сигнал который мониторит принудительную остановку программы и отправляет в контекст
	// stop() все равно остановит программу если вдруг сигнал не отправится
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// запускается в фоне и будет остановлена только после принудительной остановки программы
	url := "http://" + cfg.Address
	go metrics.Monitor(ctx, url, cfg.PoolInterval, cfg.ReportInterval)
	<-ctx.Done()
}
