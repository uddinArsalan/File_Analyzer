package main

import (
	"context"
	"file-analyzer/internals/adapters/backblaze"
	"file-analyzer/internals/adapters/cohere"
	"file-analyzer/internals/adapters/qdrant"
	"file-analyzer/internals/adapters/redis"
	"file-analyzer/internals/db/db"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	// ENV Variables
	err := godotenv.Load()
	l := log.New(os.Stdout, "DOC WORKER: ", log.LstdFlags|log.Lshortfile)
	if err != nil {
		l.Fatal("Error loading .env file", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// DB Connection
	dbClient, err := db.NewDBConnection(l)
	if err != nil {
		l.Fatal(err)
	}
	defer dbClient.CloseDB()

	// Qdrant Vector
	collection := os.Getenv("COLLECTION_NAME")
	qClient, err := qdrant.NewQdrantClient(ctx, collection)
	if err != nil {
		l.Fatal("Error Initialising Qdrant Client", err)
	}

	// Cohere Client
	cohereClient, err := cohere.NewCohereClient(ctx)
	if err != nil {
		l.Fatal("Error Initialising Cohere Client ", err)
	}

	// S3 Client
	s3Client, err := backblaze.NewS3Client(ctx)
	if err != nil {
		l.Fatal("Error Initialising Backblaze Client ", err)
	}

	rdb, err := redis.NewRedisClient(ctx)
	if err != nil {
		l.Fatal(err)
	}

	if err := rdb.CreateAndCheckStream(ctx); err != nil {
		if strings.Contains(err.Error(), "BUSYGROUP") {
			l.Println("BUSYGROUP Consumer group already exists")
		} else {
			l.Fatal(err.Error())
		}
	}

	d := NewDispatcher(ctx, 3, 12)
	d.Start(l, cohereClient, qClient, dbClient, s3Client, rdb)
	// d.StartRedisListener(ctx, l, rdb)

	<-ctx.Done()

	log.Println("Shutting down dispatcher...")
	d.Stop()
	log.Println("Graceful shutdown complete")

}
