package cohere

import (
	"context"

	cohere "github.com/cohere-ai/cohere-go/v2"
	"github.com/qdrant/go-client/qdrant"
)

type Embedder interface {
	GenerateEmbedding(ctx context.Context, text []string, inputType cohere.EmbedInputType) (*cohere.EmbedByTypeResponse, error)
	ProcessChunks(ctx context.Context, userId, docId string, chunksText []string) ([]*qdrant.PointStruct, error)
}
