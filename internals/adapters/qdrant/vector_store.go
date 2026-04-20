package qdrant

import (
	"context"
	"file-analyzer/internals/domain"

	"github.com/qdrant/go-client/qdrant"
)

type VectorStore interface {
	InsertVectorEmbeddings(ctx context.Context, vectorPoints []domain.VectorPoint) error
	SearchEmbeddingInDocument(ctx context.Context, embedding []float64, docId string) ([]*qdrant.ScoredPoint, error)
}
