package handlers

import (
	"encoding/json"
	"file-analyzer/internals/handlers/dto"
	"file-analyzer/internals/middlewares"
	"file-analyzer/internals/services"
	"file-analyzer/internals/utils"
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

func (h *UserFileHandler) FileHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")

	if err != nil {
		utils.FAIL(w, http.StatusBadRequest, "Invalid file")
		return
	}

	defer file.Close()

	userId := r.Context().Value(middlewares.UserID{}).(string)
	docId, err := h.service.UploadAndProcess(r.Context(), file, userId)

	if err != nil {
		h.l.Printf("upload failed: %v", err)
		utils.FAIL(w, http.StatusInternalServerError, "Upload failed")
		return
	}

	utils.SUCCESS(w, "File Uploaded Successfully", dto.FileResponse{
		DocID: docId,
	})
}

func (h *UserFileHandler) GenerateHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.DocRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.FAIL(w, http.StatusBadRequest, "Invalid File Details")
		return
	}
	userId := r.Context().Value(middlewares.UserID{}).(string)
	url, err := h.service.GeneratePresignedURL(r.Context(), userId, req.FileName)

	if err != nil {
		h.l.Printf("Generate failed: %v", err)
		utils.FAIL(w, http.StatusInternalServerError, "Failed to Generate Presigned URL")
		return
	}

	utils.SUCCESS(w, "Presigned URL Generated Successfully", dto.PresignedResponse{
		UploadURL: url,
	})

}
