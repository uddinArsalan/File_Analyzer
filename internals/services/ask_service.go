package services

import (
	"context"
	llm "file-analyzer/internals/adapters/cohere"
	"file-analyzer/internals/adapters/qdrant"
	"file-analyzer/internals/domain"
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

func (s *AskService) Ask(ctx context.Context, question string, docId string) (*domain.AskResponse, error) {
	embeddings, err := s.llm.GenerateEmbedding(ctx, []string{question}, domain.EmbedInputTypeSearchQuery)

	if err != nil {
		return nil, err
	}

	response, err := s.vector.SearchEmbeddingInDocument(ctx, embeddings[0], docId)
	if err != nil {
		return nil, err
	}
	docs := make([]string, len(response))
	for i, res := range response {
		docs[i] = res.Payload
	}
	reranked, err := s.llm.RerankContext(ctx, docs, question)
	if err != nil {
		return nil, err
	}

	result, err := s.llm.GenerateResponse(ctx, question, reranked)
	if err != nil {
		return nil, err
	}
	return &domain.AskResponse{
		Content:   result.Content,
		Citations: result.Citations,
	}, nil
}
