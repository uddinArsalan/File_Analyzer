package qdrant

import (
	"context"

	"github.com/qdrant/go-client/qdrant"
)

type VectorStore interface {
	InsertVectorEmbeddings(ctx context.Context, points []*qdrant.PointStruct) (*qdrant.UpdateResult, error)
	SearchEmbeddingInDocument(ctx context.Context, embedding []float64, docId string) ([]*qdrant.ScoredPoint, error)
}
