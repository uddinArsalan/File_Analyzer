package domain

import (
	"fmt"
)

type EmbeddingMetaData struct {
	Embeddings []float64
	ChunkID    string
	DocID      string
	UserID     int64
	Text       string
}

// Specifies the type of input passed to the model.
//
// - `"search_document"`: Used for embeddings stored in a vector database for search use-cases.
// - `"search_query"`: Used for embeddings of search queries run against a vector DB to find relevant documents.
// - `"classification"`: Used for embeddings passed through a text classifier.
// - `"clustering"`: Used for the embeddings run through a clustering algorithm.
// - `"image"`: Used for embeddings with image input.
type EmbedInputType string

const (
	EmbedInputTypeSearchDocument EmbedInputType = "search_document"
	EmbedInputTypeSearchQuery    EmbedInputType = "search_query"
	EmbedInputTypeClassification EmbedInputType = "classification"
	EmbedInputTypeClustering     EmbedInputType = "clustering"
	EmbedInputTypeImage          EmbedInputType = "image"
)

func NewEmbedInputTypeFromString(s string) (EmbedInputType, error) {
	switch s {
	case "search_document":
		return EmbedInputTypeSearchDocument, nil
	case "search_query":
		return EmbedInputTypeSearchQuery, nil
	case "classification":
		return EmbedInputTypeClassification, nil
	case "clustering":
		return EmbedInputTypeClustering, nil
	case "image":
		return EmbedInputTypeImage, nil
	}
	var t EmbedInputType
	return "", fmt.Errorf("%s is not a valid %T", s, t)
}
