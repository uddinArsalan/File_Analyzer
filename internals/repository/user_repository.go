package repo

import (
	"file-analyzer/internals/domain"
	"time"
)

type UserRepository interface {
	FindUserByEmail(email string) (domain.User, error)
	FindUserById(userId string) (domain.User, error)
	InsertUser(user domain.User) error
	InsertRefreshToken(tokenHash string, userID int64, ttl time.Duration) (int64, error)
	FindUserByToken(token string) (domain.RefreshToken, error)
	RevokeRefreshToken(oldTokenID int64, newTokenID int64) error
	InsertDoc(docID string, doc domain.Document) error
	UpdateDocStatus(docID string, status string) error
	DocumentExistsForUser(userID int64, docID string) error
}
