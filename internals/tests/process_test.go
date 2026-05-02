package tests

import (
	"file-analyzer/internals/domain"
	"file-analyzer/internals/utils"
	"log"
	"strings"
	"testing"
)

func TestProcessChunk(t *testing.T) {
	chunks := []domain.Chunks{
		{
			ChunkIndex: 0,
			ChunkText:  "This is the first chunk of text.",
			MetaData: map[domain.MetaDataKeys]interface{}{
				domain.UserIDKey: int64(1),
				domain.DocIDKey:  "doc-001",
			},
		},
		{
			ChunkIndex: 1,
			ChunkText:  "This is the second chunk of text.",
			MetaData: map[domain.MetaDataKeys]interface{}{
				domain.UserIDKey: int64(1),
				domain.DocIDKey:  "doc-001",
			},
		},
		{
			ChunkIndex: 2,
			ChunkText:  strings.Repeat("x", 600),
			MetaData: map[domain.MetaDataKeys]interface{}{
				domain.UserIDKey: int64(1),
				domain.DocIDKey:  "doc-001",
			},
		},
		{
			ChunkIndex: 3,
			ChunkText:  "Short chunk after a long one.",
			MetaData: map[domain.MetaDataKeys]interface{}{
				domain.UserIDKey: int64(1),
				domain.DocIDKey:  "doc-001",
			},
		},
	}
	output := utils.BatchChunksForEmbedding(chunks)
	log.Printf("Batches %+v", output)

}
