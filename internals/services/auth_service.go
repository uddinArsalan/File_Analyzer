package services

import (
	"database/sql"
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
	refreshToken := ""
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
