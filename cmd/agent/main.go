// ЦЕЛИ: написать максимально идоматичный код go
// 		 изучить библиотеку runtime
// TODO Я начал делать агента раньше времени )))))) МБ придется все переписать когда дадут теорию для написания агента

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"time"
	// "net/http"
)

type gauge struct {
	name string
	data float64
}

type counter struct {
	name string
	data int64
}

func SendData[metricData gauge | counter](m metricData, url string) {
	// логика отправки метрик на сервер
}

func Monitor(ctx context.Context, interval time.Duration) {
	ticker := time.Tick(interval * time.Second)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	for {
		select {
		case <-ticker:
			// пример данных которые можно отправить на сервер
			fmt.Println("===========- SYSTEM STATUS -===========")
			fmt.Println("Active goroutines:", runtime.NumGoroutine())
			fmt.Println("Memory usage:", m.Alloc/1024)
			fmt.Println("Memory reserved:", m.Sys/1024)
			fmt.Println("CPU allowed:", runtime.NumCPU())

			// тут вызов функции SendData с переданными ей метриками

		case <-ctx.Done():
			fmt.Println("Agent stopped")
			return
		}
	}
}

func main() {
	// для запуска программы надо использовать флаг -metrics 5 (5 - интервал для тиков)
	metricsFlag := flag.Int("metrics", 1, "collect metrics per N second")
	flag.Parse()

	// сигнал который мониторит принудительную остановку программы и отправляет в контекст
	// stop() все равно остановит программу если вдруг сигнал не отправится
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// запускается в фоне и будет остановлена только после принудительной остановки программы
	go Monitor(ctx, time.Duration(*metricsFlag))
	<-ctx.Done()
}
