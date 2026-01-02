package server

import (
	"encoding/json"
	"file-analyzer/internals/cohere"
	db "file-analyzer/internals/db/qdrant"
	"file-analyzer/internals/handlers"
	"file-analyzer/internals/middlewares"
	"file-analyzer/internals/utils"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Server struct {
	r      *chi.Mux
	Qdrant *db.QdrantClient
	Cohere *cohere.UserClient
	logger *log.Logger
}

func NewServer(r *chi.Mux, qClient *db.QdrantClient, cohereClient *cohere.UserClient, logger *log.Logger) *Server {
	s := &Server{
		r:      r,
		Qdrant: qClient,
		Cohere: cohereClient,
		logger: logger,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	userFileHandler := handlers.NewFileHandler(s.Qdrant, s.Cohere, s.logger)
	askHandler := handlers.NewAskHandler(s.Qdrant, s.Cohere,s.logger)

	// middlewares
	s.r.Use(middlewares.RateLimiter)

	// CORS
	var allowedOrigins []string
	json.Unmarshal([]byte(os.Getenv("ALLOWED_ORIGINS_JSON")), &allowedOrigins)

	s.r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: false,
	}))

	s.r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.SUCCESS(w, "welcome", nil)
	})

	s.r.Post("/ask/{docId}", askHandler.Askandler)

	s.r.Post("/upload", userFileHandler.FileHandler)
}
