package cohere

import (
	"context"
	"file-analyzer/internals/domain"
	cohere "github.com/cohere-ai/cohere-go/v2"
	"github.com/qdrant/go-client/qdrant"
)

type Embedder interface {
	GenerateEmbedding(ctx context.Context, text []string, inputType cohere.EmbedInputType) (*cohere.EmbedByTypeResponse, error)
	ProcessChunks(ctx context.Context, chunks []domain.Chunks) ([]*qdrant.PointStruct, error)
	RerankContext(ctx context.Context, documents []string, question string) ([]string, error)
	GenerateResponse(ctx context.Context, userQuestion string, documents []string) (*cohere.AssistantMessageResponse, error)
}
