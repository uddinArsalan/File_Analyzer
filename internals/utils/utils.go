package utils

import (
	"encoding/json"
	"file-analyzer/internals/domain"
	"net/http"
)

type ApiResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data,omitempty"`
	SuccessMsg string      `json:"message,omitempty"`
	Error      string      `json:"error,omitempty"`
}

type LLMPrompt struct {
	Prompt string
}

func JSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func SUCCESS(w http.ResponseWriter, status int, successMsg string, data interface{}) {
	JSON(w, status, ApiResponse{
		Success:    true,
		SuccessMsg: successMsg,
		Data:       data,
	})
}

func FAIL(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, ApiResponse{
		Success: false,
		Error:   msg,
	})
}

const MAX_CHUNK_SIZE = 500

const MAX_BATCH_SIZE = 96

func BatchChunksForEmbedding(chunks []domain.Chunks) [][]domain.Chunks {
	var batches = [][]domain.Chunks{}
	var currentBatch = []domain.Chunks{}
	for i, chunk := range chunks {
		chunkLength := len(chunk.ChunkText)
		if chunkLength > MAX_CHUNK_SIZE {
			j := 0
			for j < chunkLength {
				end := min(j+MAX_CHUNK_SIZE, chunkLength)
				currChunkText := chunk.ChunkText[j:end]
				currentBatch = append(currentBatch, domain.Chunks{
					ChunkID:   chunk.ChunkID,
					MetaData:  chunk.MetaData,
					ChunkText: currChunkText,
				})
				j += MAX_CHUNK_SIZE
			}
			j += MAX_CHUNK_SIZE
		}
		if i%MAX_BATCH_SIZE == 0 {
			batches = append(batches, currentBatch)
			currentBatch = []domain.Chunks{}
		}
	}
	return batches
}
