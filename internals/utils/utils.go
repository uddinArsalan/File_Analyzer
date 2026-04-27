package utils

import (
	"encoding/json"
	"file-analyzer/internals/domain"
	"net/http"
	"fmt"
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
	for _, chunk := range chunks {
		runes := []rune(chunk.ChunkText)
		chunkLength := len(runes)
		if chunkLength > MAX_CHUNK_SIZE {
			j := 0
			subIndex := 0
			for j < chunkLength {
				end := min(j+MAX_CHUNK_SIZE, chunkLength)
				currChunkText := string(runes[j:end])
				if len(currentBatch) == MAX_BATCH_SIZE {
					batches = append(batches, currentBatch)
					currentBatch = []domain.Chunks{}
				}
				chunkId := fmt.Sprintf("%s_%d",chunk.ChunkID,subIndex)
				currentBatch = append(currentBatch, domain.Chunks{
					ChunkID:  chunkId,
					MetaData:  chunk.MetaData,
					ChunkText: currChunkText,
				})
				subIndex++
				j += MAX_CHUNK_SIZE
			}
		} else {
			if len(currentBatch) == MAX_BATCH_SIZE {
				batches = append(batches, currentBatch)
				currentBatch = []domain.Chunks{}
			}
			currentBatch = append(currentBatch, chunk)
		}
	}
	if len(currentBatch) > 0 {
		batches = append(batches, currentBatch)
	}
	return batches
}
