package domain

type VectorPoint struct {
	Id      string
	Vectors []float32
	Payload map[string]any
}
