package services

import (
	"context"
	llm "file-analyzer/internals/adapters/cohere"
	"file-analyzer/internals/adapters/qdrant"

	cohere "github.com/cohere-ai/cohere-go/v2"
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

func (s *AskService) Ask(ctx context.Context, question string, docId string) (*cohere.AssistantMessageResponse, error) {
	resp, err := s.llm.GenerateEmbedding(ctx, []string{question}, cohere.EmbedInputTypeSearchQuery)

	if err != nil {
		return nil, err
	}

	response, err := s.vector.SearchEmbeddingInDocument(ctx, resp.Embeddings.Float[0], docId)
	if err != nil {
		return nil, err
	}
	docs := make([]string, len(response))
	for i, res := range response {
		docs[i] = res.Payload["text"].GetStringValue()
	}
	reranked, err := s.llm.RerankContext(ctx, docs, question)
	if err != nil {
		return nil, err
	}

	result, err := s.llm.GenerateResponse(ctx, question, reranked)
	if err != nil {
		return nil, err
	}
	return result, nil
}
