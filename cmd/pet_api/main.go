package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/marviel-vananaz/go-stack-backend/.gen/api"
	"github.com/marviel-vananaz/go-stack-backend/infra/sqlite"
	"github.com/marviel-vananaz/go-stack-backend/usecase/petsvc"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/api/option"
)

type Middleware func(http.Handler) http.Handler

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func withCORS() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func withFirebaseAuth(authClient *auth.Client) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			idToken := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := authClient.VerifyIDToken(r.Context(), idToken)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid token: %s", token), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userId", token.UID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func main() {
	// Initialize Firebase Admin SDK
	opt := option.WithCredentialsFile("./firebase_config.json") // Update this path
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

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

	// Apply middlewares in a more readable way
	handler := Chain(srv,
		withCORS(),
		withFirebaseAuth(authClient),
	)

	port := 8080
	fmt.Printf("Listening to port: %d \n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), handler); err != nil {
		log.Fatal(err)
	}
}
