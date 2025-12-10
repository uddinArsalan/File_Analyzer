package main

import (
	"context"
	"file-analyzer/internals/cohere"
	db "file-analyzer/internals/db/qdrant"
	"file-analyzer/internals/server"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error loading .env file")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	qClient, err := db.NewQdrantClient(ctx)
	if err != nil {
		fmt.Println("Error Initialising Qdrant Client", err)
	}

	cohereClient, err := cohere.NewCohereClient(ctx)
	if err != nil {
		fmt.Println("Error Initialising Cohere Client ", err)
	}

	exists, err := qClient.CollectionExists(ctx)
	if err != nil {
		log.Println("Error checking collection:", err)
		return
	}
	if !exists {
		err := qClient.CreateCollection(ctx)
		if err != nil {
			log.Println("Error Creating Collection", err)
			return
		}
	}

	defer qClient.Close()

	r := chi.NewRouter()

	server.NewServer(r, qClient.Qdrant, cohereClient.Cohere)

	err = http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatal("Server Exit")
	}

	fmt.Println("Server listening on port 3000")
}
