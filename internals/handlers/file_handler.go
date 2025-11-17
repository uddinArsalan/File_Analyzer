package handlers

import (
	// "file-analyzer/internals/config"
	"file-analyzer/internals/cohere"
	db "file-analyzer/internals/db/qdrant"
	"file-analyzer/internals/utils"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
)

type UserFileHandler struct {
	Qdrant *db.QdrantClient
	Cohere *cohere.UserClient
}

func NewFileHandler(qClient *db.QdrantClient, cohereClient *cohere.UserClient) *UserFileHandler {
	return &UserFileHandler{Qdrant: qClient, Cohere: cohereClient}
}

func (f *UserFileHandler) FileHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		utils.FAIL(w, http.StatusBadRequest, "Failed to read file")
        return
	}
	// generate doc id (unique) for each document
	// userId := r.Context().Value("userId").(string)
	userId := "1"
	docId := uuid.New().String()

	const MAX_CHUNKS = 96
	var (
		buff        = make([]byte, 4096)
		builder     strings.Builder
		chunkBuffer []string
	)
	for {
		n, err := file.Read(buff)
		if n > 0 {
			builder.Write(buff[:n])
			if builder.Len() >= 400 {
				chunkBuffer = append(chunkBuffer, builder.String())
				builder.Reset()
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			utils.FAIL(w, http.StatusInternalServerError, "Failed to read file")
			return
		}
	}
	for i := 0; i < len(chunkBuffer); i += MAX_CHUNKS {
		end := min(i+MAX_CHUNKS, len(chunkBuffer))
		chunks := chunkBuffer[i:end]
		points, err := f.Cohere.ProcessChunks(r.Context(), userId, docId, chunks)
		if err != nil {
			utils.FAIL(w, http.StatusInternalServerError, "Failed to process embeddings")
            return
		}
		// after it store in db (doc id , user id ,file meta info ) maybe
		res, err := f.Qdrant.InsertVectorEmbeddings(points)
		fmt.Println("Response ", res)
		if err != nil {
			utils.FAIL(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	utils.SUCCESS(w, "File Uploaded Successfully", nil)
}
