package handlers

import (
	"file-analyzer/internals/handlers/dto"
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

	userId := r.Context().Value("userId").(string)
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
