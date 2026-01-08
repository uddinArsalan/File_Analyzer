package server

import (
	"encoding/json"
	"file-analyzer/internals/adapters/cohere"
	"file-analyzer/internals/adapters/jwt"
	"file-analyzer/internals/adapters/qdrant"
	"file-analyzer/internals/handlers"
	"file-analyzer/internals/middlewares"
	repo "file-analyzer/internals/repository"
	"file-analyzer/internals/services"
	"file-analyzer/internals/utils"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Server struct {
	router *chi.Mux
	logger *log.Logger
}

func NewServer(router *chi.Mux, qdrantClient qdrant.VectorStore, embedder cohere.Embedder, userRepo repo.UserRepository, logger *log.Logger, tokenService jwt.TokenService) *Server {
	s := &Server{
		router: router,
		logger: logger,
	}
	s.routes(qdrantClient, embedder, tokenService, userRepo)
	return s
}

func (s *Server) routes(qdrantClient qdrant.VectorStore, embedder cohere.Embedder, tokenService jwt.TokenService, userRepo repo.UserRepository) {
	// services
	fileService := services.NewFileService(qdrantClient, embedder)
	askService := services.NewAskService(qdrantClient, embedder)
	authService := services.NewAuthService(userRepo, tokenService)

	userFileHandler := handlers.NewFileHandler(fileService, s.logger)
	askHandler := handlers.NewAskHandler(askService, s.logger)
	authHandler := handlers.NewAuthHandler(s.logger, authService)

	// middlewares
	s.router.Use(middlewares.RateLimiter)

	// CORS
	var allowedOrigins []string
	if err := json.Unmarshal([]byte(os.Getenv("ALLOWED_ORIGINS_JSON")), &allowedOrigins); err != nil {
		s.logger.Println("invalid ALLOWED_ORIGINS_JSON", err)
	}

	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: false,
	}))

	s.router.Get("/health", middlewares.Auth(func(w http.ResponseWriter, r *http.Request) {
		utils.SUCCESS(w, "All good", nil)
	}))

	s.router.Post("/ask/{docId}", askHandler.AskHandler)
	s.router.Post("/auth/login", authHandler.LoginHandler)
	s.router.Post("/auth/register", authHandler.RegisterHandler)

	s.router.Post("/upload", userFileHandler.FileHandler)
}
