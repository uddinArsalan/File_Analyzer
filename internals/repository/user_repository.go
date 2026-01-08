package repo

import (
	"file-analyzer/internals/domain"
)

type UserRepository interface {
	FindUserByEmail(email string) (domain.User, error)
	FindUserById(userId string) (domain.User, error)
	InsertUser(user domain.User) error
}
