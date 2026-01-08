package domain

import "time"

type Document struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	DocUrl    string    `json:"doc_url"`
	DocSize   int64     `json:"doc_size"`
	CreatedAt time.Time `json:"created_at"`
}
