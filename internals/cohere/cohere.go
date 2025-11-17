package cohere

import (
	"context"
	"fmt"
	// db "file-analyzer/internals/db/qdrant"
	"log"
	"os"

	cohere "github.com/cohere-ai/cohere-go/v2"
	"github.com/cohere-ai/cohere-go/v2/client"
	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
)

type Payload struct {
	userId  string
	docId   string
	chunkId string
}

type UserClient struct {
	Cohere *client.Client
}

func NewCohereClient(ctx context.Context) (*UserClient, error) {
	cohereClient := client.NewClient(client.WithToken(os.Getenv("CO_API_KEY")))
	return &UserClient{Cohere: cohereClient}, nil
}

func (cc *UserClient) ProcessChunks(ctx context.Context, userId, docId string, chunksText []string) ([]*qdrant.PointStruct, error) {
	resp, err := cc.Cohere.V2.Embed(
		context.TODO(),
		&cohere.V2EmbedRequest{
			Texts:          chunksText,
			Model:          "embed-v4.0",
			InputType:      cohere.EmbedInputTypeSearchDocument,
			EmbeddingTypes: []cohere.EmbeddingType{cohere.EmbeddingTypeFloat},
		},
	)
	if err != nil {
		log.Printf("Failed to generate embeddings: %v", err)
		return nil, fmt.Errorf("embedding generation failed: %w", err)
	}

	float64Emb := resp.Embeddings.Float[0]
	float32Emb := make([]float32, len(float64Emb))
	for i, v := range float64Emb {
		float32Emb[i] = float32(v)
	}

	points := make([]*qdrant.PointStruct, 0, len(resp.Embeddings.Float))
	for _, emb := range float32Emb {
		chunkId := uuid.New().String()

		point := &qdrant.PointStruct{
			Id:      qdrant.NewID(chunkId),
			Vectors: qdrant.NewVectors(emb),
			Payload: qdrant.NewValueMap(map[string]any{
				"user_id":  userId,
				"doc_id":   docId,
				"chunk_id": chunkId,
			}),
		}
		points = append(points, point)
	}
	// log.Printf("%+v", resp)
	return points, nil
}
