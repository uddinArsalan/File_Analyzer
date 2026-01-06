package main

import (
	"context"
	"file-analyzer/internals/cohere"
	"file-analyzer/internals/db/db"
	qdrant "file-analyzer/internals/db/qdrant"
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

	dbClient, err := db.NewDBConnection(l)
	if err != nil {
		l.Fatal(err)
	}
	defer dbClient.CloseDB()

	qClient, err := qdrant.NewQdrantClient(ctx)
	if err != nil {
		l.Fatal("Error Initialising Qdrant Client", err)
	}

	cohereClient, err := cohere.NewCohereClient(ctx)
	if err != nil {
		l.Fatal("Error Initialising Cohere Client ", err)
	}

	exists, err := qClient.CollectionExists(ctx)
	if err != nil {
		l.Println("Error checking collection:", err)
		return
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

	server.NewServer(r, qClient, cohereClient, l, dbClient)

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
	l.Printf("Gracefully Shutdown , Recieve Signal : %v", sign)

	ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	s.Shutdown(ctx)
}
