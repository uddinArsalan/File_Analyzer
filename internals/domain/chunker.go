package domain

type Chunks struct {
	ChunkID   string
	MetaData  map[string]interface{}
	ChunkText string
}