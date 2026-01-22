package server

import (
	"encoding/json"
	"file-analyzer/internals/adapters/backblaze"
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

func NewServer(router *chi.Mux, qdrantClient qdrant.VectorStore, embedder cohere.Embedder, userRepo repo.UserRepository, s3Client backblaze.S3Store, logger *log.Logger, tokenService jwt.TokenService) *Server {
	s := &Server{
		router: router,
		logger: logger,
	}
	s.routes(qdrantClient, embedder, s3Client, tokenService, userRepo)
	return s
}

func (s *Server) routes(qdrantClient qdrant.VectorStore, embedder cohere.Embedder, s3Client backblaze.S3Store, tokenService jwt.TokenService, userRepo repo.UserRepository) {
	// services
	fileService := services.NewFileService(qdrantClient, embedder, s3Client)
	askService := services.NewAskService(qdrantClient, embedder)
	authService := services.NewAuthService(userRepo, tokenService)

	userFileHandler := handlers.NewFileHandler(fileService, s.logger)
	askHandler := handlers.NewAskHandler(askService, s.logger)
	authHandler := handlers.NewAuthHandler(s.logger, authService)

	// CORS
	var allowedOrigins []string
	if err := json.Unmarshal([]byte(os.Getenv("ALLOWED_ORIGINS_JSON")), &allowedOrigins); err != nil {
		s.logger.Println("invalid ALLOWED_ORIGINS_JSON", err)
	}

	s.router.Group(func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   allowedOrigins,
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowCredentials: false,
		}))

		// PUBLIC ROUTES GROUP
		r.Group(func(r chi.Router) {
			r.Post("/auth/login", authHandler.LoginHandler)
			r.Post("/auth/register", authHandler.RegisterHandler)
			r.Post("/auth/refresh", authHandler.RefreshHandler)
		})

		// PRIVATE ROUTES GROUP
		r.Group(func(r chi.Router) {
			r.Use(middlewares.Auth(*authService))

			r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
				s.logger.Println(r.Context().Value(middlewares.UserID{}))
				utils.SUCCESS(w, "All good", nil)
			})

			// DOC ROUTES
			r.Post("/ask/{docId}", askHandler.AskHandler)
			r.Post("/upload", userFileHandler.FileHandler)

			// Presigned URL
			r.Post("/generate", userFileHandler.GenerateHandler)
		})

	})

}
