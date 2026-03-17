package cohere

import (
	"context"
	"file-analyzer/internals/domain"
	"fmt"
	"log"
	"os"

	cohere "github.com/cohere-ai/cohere-go/v2"
	"github.com/cohere-ai/cohere-go/v2/client"
	"github.com/qdrant/go-client/qdrant"
)

type UserClient struct {
	Cohere *client.Client
}

func NewCohereClient(ctx context.Context) (*UserClient, error) {
	cohereClient := client.NewClient(client.WithToken(os.Getenv("CO_API_KEY")))
	return &UserClient{Cohere: cohereClient}, nil
}

func (cc *UserClient) GenerateEmbedding(ctx context.Context, text []string, inputType cohere.EmbedInputType) (*cohere.EmbedByTypeResponse, error) {
	resp, err := cc.Cohere.V2.Embed(
		ctx,
		&cohere.V2EmbedRequest{
			Texts:          text,
			Model:          "embed-v4.0",
			InputType:      inputType,
			EmbeddingTypes: []cohere.EmbeddingType{cohere.EmbeddingTypeFloat},
		},
	)
	if err != nil {
		log.Printf("Failed to generate embeddings: %v", err)
		return nil, fmt.Errorf("Embedding generation failed: %w", err)
	}
	return resp, nil
}

func (cc *UserClient) ProcessChunks(ctx context.Context, chunks []domain.Chunks) ([]*qdrant.PointStruct, error) {

	texts := make([]string, 0, len(chunks))

	for _, chunk := range chunks {
		texts = append(texts, chunk.ChunkText)
	}

	resp, err := cc.GenerateEmbedding(ctx, texts, cohere.EmbedInputTypeSearchDocument)
	if err != nil {
		return nil, err
	}

	points := make([]*qdrant.PointStruct, 0, len(resp.Embeddings.Float))

	for i, float64Vectors := range resp.Embeddings.Float {
		vector := make([]float32, len(float64Vectors))

		for j, v := range float64Vectors {
			vector[j] = float32(v)
		}
		point := &qdrant.PointStruct{
			Id:      qdrant.NewID(chunks[i].ChunkID),
			Vectors: qdrant.NewVectors(vector...),
			Payload: qdrant.NewValueMap(map[string]any{
				"user_id":  chunks[i].MetaData["user_id"],
				"doc_id":   chunks[i].MetaData["doc_id"],
				"chunk_id": chunks[i].ChunkID,
				"text":     chunks[i].ChunkText,
			}),
		}
		points = append(points, point)
	}

	return points, nil
}
