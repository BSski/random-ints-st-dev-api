package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Needed for runs without Docker.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	port := "8080"
	if portFromEnv := os.Getenv("PORT"); portFromEnv != "" {
		port = portFromEnv
	}

	log.Printf("Starting up on http://localhost:%s", port)

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello World!"))
	})
	rs := randomAPI.randomAPIResource{}
	r.Mount("/random", rs.randomAPI.Routes())

	log.Fatal(http.ListenAndServe(":"+port, r))
}
