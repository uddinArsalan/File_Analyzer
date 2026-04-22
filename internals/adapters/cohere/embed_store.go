package cohere

import (
	"context"
	"file-analyzer/internals/domain"
)

type Embedder interface {
	GenerateEmbedding(ctx context.Context, text []string, inputType domain.EmbedInputType) ([][]float64, error)
	ProcessChunks(ctx context.Context, chunks []domain.Chunks) ([]domain.VectorPoint, error)
	RerankContext(ctx context.Context, documents []string, question string) ([]string, error)
	GenerateResponse(ctx context.Context, userQuestion string, documents []string) (*domain.AskResponse, error)
}
