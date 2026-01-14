package services

import (
	"context"
	llm "file-analyzer/internals/adapters/cohere"
	"file-analyzer/internals/adapters/qdrant"
	cohere "github.com/cohere-ai/cohere-go/v2"
	"time"
)

type AskService struct {
	vector qdrant.VectorStore
	llm    llm.Embedder
}

func NewAskService(vector qdrant.VectorStore, llm llm.Embedder) *AskService {
	return &AskService{
		vector: vector,
		llm:    llm,
	}
}

func (s *AskService) Ask(question string, docId string) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := s.llm.GenerateEmbedding(ctx, []string{question}, cohere.EmbedInputTypeSearchQuery)

	if err != nil {
		return nil, err
	}

	embed := resp.Embeddings.Float[0]

	response, err := s.vector.SearchEmbedInDocument(ctx, embed, docId)
	if err != nil {
		return nil, err
	}
	//now need to create a context and send to llm to answer
	return response, nil
}
