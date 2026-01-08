package main

import (
	"context"
	"file-analyzer/internals/adapters/cohere"
	"file-analyzer/internals/adapters/jwt"
	"file-analyzer/internals/adapters/qdrant"
	"file-analyzer/internals/db/db"
	"file-analyzer/internals/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	l := log.New(os.Stdout, "DOC API: ", log.LstdFlags|log.Lshortfile)
	if err != nil {
		l.Fatal("Error loading .env file", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	secret := os.Getenv("JWT_SECRET_KEY")
	collection := os.Getenv("COLLECTION_NAME")

	dbClient, err := db.NewDBConnection(l)
	if err != nil {
		l.Fatal(err)
	}
	defer dbClient.CloseDB()

	qClient, err := qdrant.NewQdrantClient(ctx, collection)
	if err != nil {
		l.Fatal("Error Initialising Qdrant Client", err)
	}

	cohereClient, err := cohere.NewCohereClient(ctx)
	if err != nil {
		l.Fatal("Error Initialising Cohere Client ", err)
	}

	exists, err := qClient.CollectionExists(ctx)
	if err != nil {
		l.Fatal("Error checking collection:", err)
	}
	if !exists {
		err := qClient.CreateCollection(ctx)
		if err != nil {
			l.Fatal("Error Creating Collection", err)
		}
	}

	var payloadFields = []string{"doc_id", "user_id", "org_id"}
	for _, fieldName := range payloadFields {
		err := qClient.EnsurePayLoadIndex(ctx, fieldName)
		if err != nil {
			l.Printf("Index create skipped for %s: %v\n", fieldName, err)
		}
	}

	defer qClient.Close()

	r := chi.NewRouter()

	tokenService := jwt.NewJwtService(secret)

	server.NewServer(r, qClient, cohereClient, dbClient, l, tokenService)

	s := &http.Server{
		Addr:         ":3000",
		Handler:      r,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  20 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatalf("Server Exit %v", err)
		}
	}()

	l.Println("Server listening on port 3000")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)
	signal.Notify(sigChan, os.Interrupt)

	sign := <-sigChan
	l.Printf("Gracefully Shutdown , Received Signal : %v", sign)

	ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	s.Shutdown(ctx)
}
