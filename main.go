package main

import (
	"log"
	"net/http"
	"os"

	"github.com/BSski/RandomIntsStDevAPI/randomintsstdev"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Get "github.com/joho/godotenv" and uncomment this code if you want to run the app without Docker.
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file", err)
	// }

	port := "8080"
	if portFromEnv := os.Getenv("PORT"); portFromEnv != "" {
		port = portFromEnv
	}

	log.Printf("Starting up on http://localhost:%s", port)

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello World! You're probably looking for '/random/mean?requests=2&length=2' endpoint."))
	})
	r.Mount("/random", randomintsstdev.RandomAPIResource{}.Routes())

	log.Fatal(http.ListenAndServe(":"+port, r))
}
