package handlers

import (
	"context"
	"encoding/json"
	"file-analyzer/internals/domain"
	"file-analyzer/internals/utils"
	"log"
	"net/http"
	"time"

	cohere "github.com/cohere-ai/cohere-go/v2"
	"github.com/go-chi/chi/v5"
)

type AskHandler struct {
	Repo domain.DocumentRepository
	LLM  domain.EmbeddingService
	l    *log.Logger
}

func NewAskHandler(repo domain.DocumentRepository, llm domain.EmbeddingService, l *log.Logger) *AskHandler {
	return &AskHandler{
		Repo: repo,
		LLM:  llm,
		l:    l,
	}
}

type UserQuestion struct {
	Question string
}

func (cc *AskHandler) Askandler(w http.ResponseWriter, r *http.Request) {
	docId := chi.URLParam(r, "docId")
	var q UserQuestion
	err := json.NewDecoder(r.Body).Decode(&q)

	defer r.Body.Close()

	if err != nil {
		cc.l.Printf("Error Reading Body %v ", err)
		utils.FAIL(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cc.l.Println(docId)
	cc.l.Println(q.Question)
	resp, err := cc.LLM.GenerateEmbedding(ctx, []string{q.Question}, cohere.EmbedInputTypeSearchQuery)
	if err != nil {
		cc.l.Printf("Embedding generation failed %v ", err)
		utils.FAIL(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	// cc.Repo

	embed := resp.Embeddings.Float[0]
	cc.l.Println(resp)
	response, err := cc.Repo.SearchEmbedInDocument(r.Context(), embed, docId)
	if err != nil {
		cc.l.Printf("Embedding Search failed %v ", err)
		utils.FAIL(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	//now need to create a context and send to llm to answer
	utils.SUCCESS(w, "Ask Successfully", response)
}
