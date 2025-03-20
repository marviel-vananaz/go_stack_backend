package main

import (
	"log"
	"net/http"

	"github.com/marviel-vananaz/go-stack-backend/internal/oas"
)

func main() {
	// Create service instance.
	service := &petsService{
		pets: map[int64]oas.Pet{},
	}
	// Create generated server.
	srv, err := oas.NewServer(service)
	if err != nil {
		log.Fatal(err)
	}
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatal(err)
	}
}
