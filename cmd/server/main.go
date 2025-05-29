package main

import (
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/database"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/router"
	"log"
	"net/http"
)

func main() {
	storage := database.NewMemStorage()
	r := router.New(storage)

	log.Fatal(http.ListenAndServe(":8080", r))
}
