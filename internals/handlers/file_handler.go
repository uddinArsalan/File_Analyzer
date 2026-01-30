package handlers

import (
	"encoding/json"
	"file-analyzer/internals/handlers/dto"
	"file-analyzer/internals/middlewares"
	"file-analyzer/internals/queue"
	"file-analyzer/internals/services"
	"file-analyzer/internals/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
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

func (h *UserFileHandler) ProcessFile(w http.ResponseWriter, r *http.Request){
	docID := r.URL.Query().Get("doc_id")
	userID := r.Context().Value(middlewares.UserID{}).(int64)
	// head request to object storage to check if file is uploaded
	isExists,err := h.service.CheckExistence(docID,userID)
	if err != nil {
		utils.FAIL(w,http.StatusNotFound,"Unable to verify file status die to storage issue.")
		return
	}
	if !isExists{
		utils.FAIL(w,http.StatusNotFound,"File has not been uploaded yet.")
		return
	}
	d := queue.NewDispatcher(3,12)
	objectKey := fmt.Sprintf("documents/%v/%v", userID, docID)
	job := queue.Job{
		ID : uuid.New().String(),
		ObjectKey: objectKey,
		UserID: userID,
		DocID: docID,
	}
	d.Submit(job)
	utils.SUCCESS(w,"File is Processing",nil)
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

	utils.SUCCESS(w, "Presigned URL Generated Successfully", dto.PresignedResponse{
		DocID:     docId,
		UploadURL: url,
	})
}
