package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenService interface {
	GenerateJWT(userId int64, ttl time.Duration) (string, error)
	VerifyToken(tokenString string) (*jwt.Token, error)
}
