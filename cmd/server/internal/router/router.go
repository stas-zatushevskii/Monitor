package router

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/database"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/gzip"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/handlers"
)

func New(storage *database.MemStorage, db *sql.DB) *chi.Mux {
	router := chi.NewRouter()
	router.Use(gzip.GzipMiddleware)

	router.Get("/", handlers.GetAllAgentHandlers(storage))
	router.Post("/update/", handlers.UpdateJSONHandler(storage))
	router.Post("/value/", handlers.ValueJSONHandler(storage))
	router.Post("/update/{type}/{name}/{data}", handlers.UpdateURLHandler(storage))
	router.Get("/value/{type}/{name}", handlers.ValueURLHandler(storage))
	router.Get("/ping", handlers.Ping(db))
	return router
}
