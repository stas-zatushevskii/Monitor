package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/database"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/router"
	"log"
	"net/http"
)

func main() {
	storage := database.NewMemStorage()
	r := router.New(storage)
	ParseFlags()

	log.Fatal(run(r))
}

func run(r *chi.Mux) error {
	fmt.Println("Running server on", address)
	return http.ListenAndServe(address, r)
}
