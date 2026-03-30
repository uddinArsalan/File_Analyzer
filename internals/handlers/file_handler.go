package handlers

import (
	"encoding/json"
	"file-analyzer/internals/handlers/dto"
	"file-analyzer/internals/middlewares"
	"file-analyzer/internals/services"
	"file-analyzer/internals/utils"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type UserFileHandler struct {
	service *services.FileService
	l       *log.Logger
}

func NewFileHandler(service *services.FileService, l *log.Logger) *UserFileHandler {
	return &UserFileHandler{
		service: service,
		l:       l,
	}
}

func (h *UserFileHandler) CheckExistenceAndProcessFile(w http.ResponseWriter, r *http.Request) {
	docID := r.URL.Query().Get("doc_id")
	userID := r.Context().Value(middlewares.UserID{}).(int64)
	err := h.service.CheckExistence(r.Context(), userID, docID)
	if err != nil {
		utils.FAIL(w, http.StatusNotFound, err.Error())
		return
	}
	utils.SUCCESS(w, http.StatusAccepted, "File is Processing", nil)
}

func (h *UserFileHandler) GenerateHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.DocRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		utils.FAIL(w, http.StatusBadRequest, "Invalid File Details")
		return
	}

	userID := r.Context().Value(middlewares.UserID{}).(int64)
	docID := uuid.New().String()
	url, err := h.service.GeneratePresignedURL(r.Context(), userID, docID, req)

	if err != nil {
		h.l.Printf("Generate failed: %v", err)
		utils.FAIL(w, http.StatusInternalServerError, "Failed to Generate Presigned URL")
		return
	}

	utils.SUCCESS(w, http.StatusOK, "Presigned URL Generated Successfully", dto.PresignedResponse{
		DocID:     docID,
		UploadURL: url,
	})
}

func (h *UserFileHandler) SSEHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserID{}).(int64)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		utils.FAIL(w, http.StatusInternalServerError, "Streaming not Supported")
		return
	}

	h.service.Stream(r.Context(), userID, flusher, w)
}
