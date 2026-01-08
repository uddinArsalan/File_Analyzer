package jwt

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret []byte
}

func NewJwtService(secret string) *JWTService {
	return &JWTService{
		secret: []byte(secret),
	}
}

func (j *JWTService) GenerateJWT(userId int64, ttl time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userId, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
	}
	s := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := s.SignedString(j.secret)
	if err != nil {
		return "", err
	}
	return token, nil
}
