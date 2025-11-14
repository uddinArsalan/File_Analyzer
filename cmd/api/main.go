package main

import (
	"encoding/json"
	"file-analyzer/internals/handlers"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error loading .env file")
	}
	r := chi.NewRouter()

	// CORS
	var allowedOrigins []string
	json.Unmarshal([]byte(os.Getenv("ALLOWED_ORIGINS_JSON")), &allowedOrigins)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: false,
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Post("/upload", handlers.FileHandler)

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Server Exit")
	}
	fmt.Println("Server listening on port 8080")
}
