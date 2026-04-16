package handlers

import (
	"encoding/json"
	"file-analyzer/internals/services"
	"file-analyzer/internals/handlers/dto"
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

func (cc *AskHandler) AskHandler(w http.ResponseWriter, r *http.Request) {
	docId := chi.URLParam(r, "docId")

	var q dto.UserQuestion
	err := json.NewDecoder(r.Body).Decode(&q)

	defer r.Body.Close()

	if err != nil {
		cc.l.Printf("Error Reading Body %v ", err)
		utils.FAIL(w, http.StatusBadRequest, "Invalid request")
		return
	}

	response, err := cc.service.Ask(r.Context(), q.Question, docId)

	utils.SUCCESS(w, http.StatusOK, "Ask Successfully", dto.AskResponse{
		Content:   response.Content,
		Citations: response.Citations,
	})
}
