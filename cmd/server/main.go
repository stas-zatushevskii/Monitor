package main

import (
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/database"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/handlers"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	storage := database.NewMemStorage()

	mux.HandleFunc("/update/", handlers.AgentHandler(storage))

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
