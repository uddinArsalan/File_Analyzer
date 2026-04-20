package tests

import (
	"file-analyzer/internals/utils"
	"log"
	"testing"
	"strings"
	"file-analyzer/internals/domain"
)

func TestProcessChunk(t *testing.T){
	chunks := []domain.Chunks{
    {
        ChunkID:  "a1b2c3d4-0001-0001-0001-000000000001",
        ChunkText: "This is the first chunk of text.",
        MetaData: map[domain.MetaDataKeys]interface{}{
            domain.UserIDKey: int64(1),
            domain.DocIDKey:  "doc-001",
        },
    },
    {
        ChunkID:  "a1b2c3d4-0002-0002-0002-000000000002",
        ChunkText: "This is the second chunk of text.",
        MetaData: map[domain.MetaDataKeys]interface{}{
            domain.UserIDKey: int64(1),
            domain.DocIDKey:  "doc-001",
        },
    },
    {
        ChunkID:  "a1b2c3d4-0003-0003-0003-000000000003",
        ChunkText: strings.Repeat("x", 600), 
        MetaData: map[domain.MetaDataKeys]interface{}{
            domain.UserIDKey: int64(1),
            domain.DocIDKey:  "doc-001",
        },
    },
    {
        ChunkID:  "a1b2c3d4-0004-0004-0004-000000000004",
        ChunkText: "Short chunk after a long one.",
        MetaData: map[domain.MetaDataKeys]interface{}{
            domain.UserIDKey: int64(1),
            domain.DocIDKey:  "doc-001",
        },
    },
}
	output := utils.BatchChunksForEmbedding(chunks)
	log.Printf("Batches %+v",output)
	
}