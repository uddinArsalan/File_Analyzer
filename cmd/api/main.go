package main

import (
	"file-analyzer/internals/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Post("/upload", handlers.FileHandler)

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Server Exit")
	}
}
