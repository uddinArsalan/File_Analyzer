package qdrant

import (
	"file-analyzer/internals/domain"

	"github.com/qdrant/go-client/qdrant"
)

func ToVectorSearchResult(qdrantResult []*qdrant.ScoredPoint) []*domain.VectorSearchResult {
	searchResult := make([]*domain.VectorSearchResult, len(qdrantResult))
	for i, res := range qdrantResult {
		searchResult[i] = &domain.VectorSearchResult{
			Payload: res.Payload["text"].GetStringValue(),
		}
	}
	return searchResult
}
