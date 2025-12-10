package server

import (
	"encoding/json"
	"file-analyzer/internals/cohere"
	db "file-analyzer/internals/db/qdrant"
	"file-analyzer/internals/handlers"
	"file-analyzer/internals/middlewares"
	"file-analyzer/internals/utils"
	"net/http"
	"os"

	"github.com/cohere-ai/cohere-go/v2/client"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/qdrant/go-client/qdrant"
)

type Server struct {
	r      *chi.Mux
	Qdrant *qdrant.Client
	Cohere *client.Client
}

func NewServer(r *chi.Mux, qClient *qdrant.Client, cohereClient *client.Client) *Server {
	s := &Server{
		r:      r,
		Qdrant: qClient,
		Cohere: cohereClient,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	userFileHandler := handlers.NewFileHandler(&db.QdrantClient{Qdrant: s.Qdrant}, &cohere.UserClient{Cohere: s.Cohere})

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

	s.r.Post("/upload", userFileHandler.FileHandler)
}
