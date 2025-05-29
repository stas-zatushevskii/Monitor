package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/database"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/handlers"
)

func New(storage *database.MemStorage) *chi.Mux {
	router := chi.NewRouter()

	// POST http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	// GET http://<АДРЕС_СЕРВЕРА>/value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
	router.Get("/", handlers.GetAllAgentHandlers(storage))
	router.Post("/update/{type}/{name}/{data}", handlers.UpdateAgentHandler(storage))
	router.Get("/value/{type}/{name}", handlers.ValueAgentHandler(storage))
	return router
}
