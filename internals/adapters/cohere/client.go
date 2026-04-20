package cohere

import (
	"context"
	"file-analyzer/internals/domain"
	"file-analyzer/internals/utils"
	"fmt"
	"log"
	"os"

	cohere "github.com/cohere-ai/cohere-go/v2"
	"github.com/cohere-ai/cohere-go/v2/client"
)

type UserClient struct {
	Cohere *client.Client
}

func NewCohereClient(ctx context.Context) (*UserClient, error) {
	cohereClient := client.NewClient(client.WithToken(os.Getenv("CO_API_KEY")))
	return &UserClient{Cohere: cohereClient}, nil
}

func (cc *UserClient) GenerateEmbedding(ctx context.Context, text []string, inputType domain.EmbedInputType) ([][]float64, error) {
	resp, err := cc.Cohere.V2.Embed(
		ctx,
		&cohere.V2EmbedRequest{
			Texts:          text,
			Model:          "embed-v4.0",
			InputType:      cohere.EmbedInputType(inputType),
			EmbeddingTypes: []cohere.EmbeddingType{cohere.EmbeddingTypeFloat},
		},
	)
	if err != nil {
		log.Printf("Failed to generate embeddings: %v", err)
		return nil, fmt.Errorf("Embedding generation failed: %w", err)
	}
	return resp.Embeddings.Float, nil
}

func (cc *UserClient) ProcessChunks(ctx context.Context, chunks []domain.Chunks) ([]domain.VectorPoint, error) {

	chunkBatches := utils.BatchChunksForEmbedding(chunks)

	accumulatedEmbeddings := make([]domain.EmbeddingMetaData, 0, len(chunks))

	for _, batch := range chunkBatches {
		texts := make([]string, len(batch))
		for i, chunk := range batch {
			texts[i] = chunk.ChunkText
		}
		log.Printf("Chunk Batch Length %d",len(texts))
		embeddings, err := cc.GenerateEmbedding(ctx, texts, domain.EmbedInputTypeSearchDocument)
		if err != nil {
			return nil, err
		}
		for i, embed := range embeddings{
			accumulatedEmbeddings = append(accumulatedEmbeddings, domain.EmbeddingMetaData{
				Embeddings: embed,
				ChunkID: batch[i].ChunkID,
				DocID:   batch[i].MetaData[domain.DocIDKey].(string),
				UserID:  batch[i].MetaData[domain.UserIDKey].(int64),
				Text:    batch[i].ChunkText,
			})
		} 
	}

	// points := make([]*qdrant.PointStruct, 0, len(accumulatedEmbeddings))
	points := make([]domain.VectorPoint, 0, len(accumulatedEmbeddings))

	for _, embed := range accumulatedEmbeddings {
		vector := make([]float32, len(embed.Embeddings))

		for j, v := range embed.Embeddings {
			vector[j] = float32(v)
		}
		// point := &qdrant.PointStruct{
		// 	Id:      qdrant.NewID(embed.ChunkID),
		// 	Vectors: qdrant.NewVectors(vector...),
		// 	Payload: qdrant.NewValueMap(map[string]any{
		// 		"user_id":  embed.UserID,
		// 		"doc_id":   embed.DocID,
		// 		"chunk_id": embed.ChunkID,
		// 		"text":     embed.Text,
		// 	}),
		// }
		point := domain.VectorPoint{
			Id:     embed.ChunkID,
			Vectors: vector,
			Payload: map[string]any{
				"user_id":  embed.UserID,
				"doc_id":   embed.DocID,
				"chunk_id": embed.ChunkID,
				"text":     embed.Text,
			},
		}
		points = append(points, point)
	}

	return points, nil
}

func (cc *UserClient) GenerateResponse(ctx context.Context, userQuestion string, documents []string) (*cohere.AssistantMessageResponse, error) {

	docs := make([]*cohere.V2ChatRequestDocumentsItem, len(documents))
	for i, text := range documents {
		docs[i] = &cohere.V2ChatRequestDocumentsItem{
			String: text,
		}
	}

	resp, err := cc.Cohere.V2.Chat(ctx, &cohere.V2ChatRequest{
		Model: "command-a-03-2025",
		Messages: cohere.ChatMessages{
			&cohere.ChatMessageV2{
				Role: "user",
				User: &cohere.UserMessageV2{
					Content: &cohere.UserMessageV2Content{
						String: userQuestion,
					},
				},
			},
		},
		Documents: docs,
	})
	if err != nil {
		log.Printf("Failed to generate response: %v", err)
		return nil, err
	}
	log.Printf("Generated response: %+v", resp)
	return resp.Message, nil
}

func (cc *UserClient) RerankContext(ctx context.Context, documents []string, question string) ([]string, error) {
	response, err := cc.Cohere.V2.Rerank(ctx, &cohere.V2RerankRequest{
		Model:     "rerank-v4.0-pro",
		Query:     question,
		Documents: documents,
		TopN:      cohere.Int(3),
	})
	if err != nil {
		log.Printf("Failed to rerank context: %v", err)
		return nil, err
	}
	log.Printf("Reranked context: %+v", response)

	reranked := make([]string, len(response.Results))
	for i, result := range response.Results {
		reranked[i] = documents[result.Index]
	}
	return reranked, nil
}
