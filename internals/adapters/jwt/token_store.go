package jwt

import "time"

type TokenService interface {
	GenerateJWT(userId int64, ttl time.Duration) (string, error)
}
