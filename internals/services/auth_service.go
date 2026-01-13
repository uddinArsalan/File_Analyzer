package services

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"file-analyzer/internals/adapters/jwt"
	"file-analyzer/internals/domain"
	repo "file-analyzer/internals/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users repo.UserRepository
	jwt   jwt.TokenService
}

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

func NewAuthService(userRepo repo.UserRepository, jwt jwt.TokenService) *AuthService {
	return &AuthService{
		users: userRepo,
		jwt:   jwt,
	}
}

var refreshExpiry = 7 * 24 * time.Hour

func (s *AuthService) Login(email, password string) (*AuthTokens, error) {
	user, err := s.users.FindUserByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			return &AuthTokens{}, ErrInvalidCredentials
		}
		return &AuthTokens{}, err
	}
	// check for password
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
	if err != nil {
		return &AuthTokens{}, ErrInvalidCredentials
	}
	accessToken, err := s.jwt.GenerateJWT(user.UserID, 5*time.Minute)
	if err != nil {
		return &AuthTokens{}, err
	}
	refreshToken, err := generateRefreshToken()
	if err != nil {
		return &AuthTokens{}, err
	}
	tokenHash := generateHash(refreshToken)
	s.users.InsertRefreshToken(tokenHash, user.UserID, refreshExpiry)
	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Register(name, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.users.InsertUser(domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: hash,
	})
}

func (s *AuthService) Refresh(incomingToken string) (*AuthTokens, error) {
	tokenHash := generateHash(incomingToken)
	oldToken, err := s.users.FindUserByToken(tokenHash)
	if err != nil {
		return &AuthTokens{}, err
	}
	if time.Now().After(oldToken.ExpiresAt) || oldToken.RevokedAt != nil {
		return &AuthTokens{}, ErrSessionExpired
	}
	newRefreshToken, err := generateRefreshToken()
	if err != nil {
		return &AuthTokens{}, err
	}
	newTokenHash := generateHash(newRefreshToken)
	newTokenId, err := s.users.InsertRefreshToken(newTokenHash, oldToken.UserID, refreshExpiry)
	if err != nil {
		return &AuthTokens{}, err
	}
	accessToken, err := s.jwt.GenerateJWT(oldToken.UserID, 5*time.Minute)
	if err != nil {
		return &AuthTokens{}, err
	}
	err = s.users.RevokeRefreshToken(oldToken.ID, newTokenId)
	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func generateHash(refreshToken string) string {
	d := sha256.Sum256([]byte(refreshToken))
	return hex.EncodeToString(d[:])
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
