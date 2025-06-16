package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/database"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/handlers"
)

func New(storage *database.MemStorage) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", handlers.GetAllAgentHandlers(storage))
	router.Post("/update", handlers.UpdateJSONHandler(storage))
	router.Post("/value", handlers.ValueJSONHandler(storage))
	router.Post("/update/{type}/{name}/{data}", handlers.UpdateURLHandler(storage))
	router.Get("/value/{type}/{name}", handlers.ValueURLHandler(storage))
	return router
}
