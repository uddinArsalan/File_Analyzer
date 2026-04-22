package qdrant

import (
	"context"
	"file-analyzer/internals/domain"
)

type VectorStore interface {
	InsertVectorEmbeddings(ctx context.Context, vectorPoints []domain.VectorPoint) error
	SearchEmbeddingInDocument(ctx context.Context, embedding []float64, docId string) ([]*domain.VectorSearchResult, error)
}
