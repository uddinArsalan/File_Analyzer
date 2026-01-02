package domain

import (
	"context"

	cohere "github.com/cohere-ai/cohere-go/v2"

	"github.com/qdrant/go-client/qdrant"
)

type DocumentRepository interface {
	InsertVectorEmbeddings(points []*qdrant.PointStruct) (*qdrant.UpdateResult, error)
	SearchEmbedInDocument(ctx context.Context, embedding []float64,docId string) ([]*qdrant.ScoredPoint, error)
}

type EmbeddingService interface {
	GenerateEmbedding(ctx context.Context, text []string, inputType cohere.EmbedInputType) (*cohere.EmbedByTypeResponse, error)
	ProcessChunks(ctx context.Context, userId, docId string, chunksText []string) ([]*qdrant.PointStruct, error)
}
