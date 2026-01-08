package domain

import "time"

type User struct {
	UserID       int64     `json:"user_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash []byte    `json:"password_hash"`
	ProfileURL   string    `json:"profile_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
