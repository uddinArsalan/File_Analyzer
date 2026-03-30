package chunker

import (
	"file-analyzer/internals/domain"
	"strings"

	"github.com/google/uuid"
)

type Chunker struct {
	DocID   string
	UserID  int64
	Content string
}

func NewChunker(text string, docID string, userID int64) *Chunker {
	return &Chunker{
		Content: text,
		DocID:   docID,
		UserID:  userID,
	}
}

func (c *Chunker) Chunk() []domain.Chunks {
	rawChunks := strings.Split(c.Content, "\n\n")
	chunks := make([]domain.Chunks, len(rawChunks))
	for i, chunkText := range rawChunks {
		chunks[i].ChunkID = uuid.NewString()
		chunks[i].ChunkText = chunkText
		chunks[i].MetaData = make(map[string]interface{})
		chunks[i].MetaData["user_id"] = c.UserID
		chunks[i].MetaData["doc_id"] = c.DocID
	}
	return chunks
}
