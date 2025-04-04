package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/marviel-vananaz/go-stack-backend/.gen/api"
	"github.com/marviel-vananaz/go-stack-backend/infra/sqlite"
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
	srv, err := api.NewServer(service)
	if err != nil {
		log.Fatal(err)
	}

	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		srv.ServeHTTP(w, r)
	})

	port := 8080
	fmt.Printf("Listening to port: %d \n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), corsHandler); err != nil {
		log.Fatal(err)
	}
}
