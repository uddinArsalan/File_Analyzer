package cohere

import (
	"context"
	"fmt"
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

func (cc *UserClient) GenerateEmbedding(ctx context.Context, text []string, inputType cohere.EmbedInputType) (*cohere.EmbedByTypeResponse, error) {
	resp, err := cc.Cohere.V2.Embed(
		context.TODO(),
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

func (cc *UserClient) ProcessChunks(ctx context.Context, userId int64, docId string, chunksText []string) ([]*qdrant.PointStruct, error) {

	// fmt.Printf("Chunks Text int ProcessChunks %v", chunksText)
	resp, err := cc.GenerateEmbedding(ctx, chunksText, cohere.EmbedInputTypeSearchDocument)
	if err != nil {
		return nil, err
	}

	points := make([]*qdrant.PointStruct, 0, len(resp.Embeddings.Float))

	for i, float64Vectors := range resp.Embeddings.Float {
		vector := make([]float32, len(float64Vectors))
		for i, v := range float64Vectors {
			vector[i] = float32(v)
		}
		// fmt.Printf("Vector i %v in ProcessChunks %v", i, vector)
		// fmt.Printf("chunksText[i] %v in ProcessChunks", chunksText[i])
		chunkId := uuid.New().String()
		point := &qdrant.PointStruct{
			Id:      qdrant.NewID(chunkId),
			Vectors: qdrant.NewVectors(vector...),
			Payload: qdrant.NewValueMap(map[string]any{
				"user_id":  userId,
				"doc_id":   docId,
				"chunk_id": chunkId,
				"text":     chunksText[i],
			}),
		}
		points = append(points, point)
	}

	return points, nil
}
