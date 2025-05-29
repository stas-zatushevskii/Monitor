package main

import (
	"context"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/metrics"
	"os"
	"os/signal"
)

func main() {
	// сигнал который мониторит принудительную остановку программы и отправляет в контекст
	// stop() все равно остановит программу если вдруг сигнал не отправится
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// запускается в фоне и будет остановлена только после принудительной остановки программы
	ParseFlags()
	url := "http://" + address
	go metrics.Monitor(ctx, url, pollIntervalFlag, reportIntervalFlag)
	<-ctx.Done()
}
