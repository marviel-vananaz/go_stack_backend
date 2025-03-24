package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/marviel-vananaz/go-stack-backend/infra/sqlite"
	"github.com/marviel-vananaz/go-stack-backend/internal/oas"
	"github.com/marviel-vananaz/go-stack-backend/usecase/petsvc"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Create service instance.
	db, err := sql.Open("sqlite3", "../database/database.db")
	if err != nil {
		panic(err)
	}
	repo := sqlite.NewPetRepo(db)
	service := petsvc.NewService(&repo)
	srv, err := oas.NewServer(service)
	if err != nil {
		log.Fatal(err)
	}
	port := 8080
	fmt.Printf("Listening to port: %d \n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), srv); err != nil {
		log.Fatal(err)
	}
}
