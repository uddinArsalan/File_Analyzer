package handlers

import (
	"file-analyzer/internals/services"
	"file-analyzer/internals/utils"
	"log"
	"net/http"
	"strconv"
)

type TempHandler struct {
	l       *log.Logger
	service *services.TempService
}

func NewTempHandler(l *log.Logger, service *services.TempService) *TempHandler {
	return &TempHandler{
		l,
		service,
	}
}

func (h *TempHandler) Add(w http.ResponseWriter, r *http.Request) {
	i := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(i, 10, 64)
	h.l.Printf("i %d", id)
	if err != nil {
		utils.FAIL(w, http.StatusBadRequest, "Bad Request")
	}
	h.service.AddTemporaryJobs(r.Context(), id)
}
