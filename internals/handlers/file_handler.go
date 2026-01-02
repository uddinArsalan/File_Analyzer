package handlers

import (
	"file-analyzer/internals/domain"
	"file-analyzer/internals/utils"
	// "fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type UserFileHandler struct {
	Repo domain.DocumentRepository
	LLM  domain.EmbeddingService
	l    *log.Logger
}

func NewFileHandler(repo domain.DocumentRepository, llm domain.EmbeddingService, l *log.Logger) *UserFileHandler {
	return &UserFileHandler{
		Repo: repo,
		LLM:  llm,
		l:    l,
	}
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
	if builder.Len() > 0 {
		chunkBuffer = append(chunkBuffer, builder.String())
		builder.Reset()
	}
	// fmt.Printf("chunkBuffer %v in FileHandler", chunkBuffer)
	for i := 0; i < len(chunkBuffer); i += MAX_CHUNKS {
		end := min(i+MAX_CHUNKS, len(chunkBuffer))
		chunks := chunkBuffer[i:end]
		// fmt.Printf("chunks passed in ProcessChunks from FileHandler %v", chunks)
		points, err := f.LLM.ProcessChunks(r.Context(), userId, docId, chunks)
		if err != nil {
			utils.FAIL(w, http.StatusInternalServerError, "Failed to process embeddings")
			return
		}
		// after it store in db (doc id , user id ,file meta info ) maybe
		res, err := f.Repo.InsertVectorEmbeddings(points)
		f.l.Println("Response ", res)
		if err != nil {
			utils.FAIL(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	utils.SUCCESS(w, "File Uploaded Successfully", docId)
}
