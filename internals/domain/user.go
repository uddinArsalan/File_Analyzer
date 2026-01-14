package domain

import "time"

type User struct {
	UserID       int64     
	Name         string    
	Email        string    
	PasswordHash []byte    
	ProfileURL   string    
	CreatedAt    time.Time 
	UpdatedAt    time.Time 
}
