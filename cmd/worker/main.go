package worker

import (
	"context"
	"file-analyzer/internals/adapters/backblaze"
	"file-analyzer/internals/adapters/cohere"
	"file-analyzer/internals/adapters/qdrant"
	"file-analyzer/internals/db/db"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main(){
	// ENV Variables
	err := godotenv.Load()
	l := log.New(os.Stdout, "DOC WORKER: ", log.LstdFlags|log.Lshortfile)
	if err != nil {
		l.Fatal("Error loading .env file", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

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

	l.Println(qClient,cohereClient,s3Client)

}