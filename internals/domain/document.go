package domain

import "time"

type Document struct {
	ID        int64
	DocID     string
	UserID    int64
	Name      string
	ObjectKey string
	Status    string
	Mime_Type string
	DocSize   int64
	CreatedAt time.Time
}
