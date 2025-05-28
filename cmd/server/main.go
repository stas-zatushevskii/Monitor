package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/database"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/handlers"
	"log"
	"net/http"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.AllowContentType("text/plain"))
	storage := database.NewMemStorage()

	// POST http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	// GET http://<АДРЕС_СЕРВЕРА>/value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
	router.Get("/", handlers.GetAllAgentHandlers(storage))
	router.Post("/update/{type}/{name}/{data}", handlers.UpdateAgentHandler(storage))
	router.Get("/value/{type}/{name}", handlers.ValueAgentHandler(storage))

	log.Fatal(http.ListenAndServe(":8080", router))
}
