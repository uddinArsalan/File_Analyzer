package domain

type VectorPoint struct {
	Id      uint64
	Vectors []float32
	Payload map[string]any
}

type VectorSearchResult struct {
	Payload string
}
