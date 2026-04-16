package domain

type EmbeddingMetaData struct {
	Embeddings []float64 
	ChunkID    string     
	DocID      string     
	UserID     int64     
	Text       string     
}