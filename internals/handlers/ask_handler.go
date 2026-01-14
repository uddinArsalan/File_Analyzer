package handlers

import (
	"encoding/json"
	"file-analyzer/internals/services"
	"file-analyzer/internals/utils"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type AskHandler struct {
	service *services.AskService
	l       *log.Logger
}

func NewAskHandler(service *services.AskService, l *log.Logger) *AskHandler {
	return &AskHandler{
		service: service,
		l:       l,
	}
}

type UserQuestion struct {
	Question string
}

func (cc *AskHandler) AskHandler(w http.ResponseWriter, r *http.Request) {
	docId := chi.URLParam(r, "docId")

	var q UserQuestion
	err := json.NewDecoder(r.Body).Decode(&q)

	defer r.Body.Close()

	if err != nil {
		cc.l.Printf("Error Reading Body %v ", err)
		utils.FAIL(w, http.StatusBadRequest, "Invalid request")
		return
	}

	response, err := cc.service.Ask(q.Question, docId)

	utils.SUCCESS(w, "Ask Successfully", response)
}
