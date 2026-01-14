package domain

import "time"

type Document struct {
	ID        int64
	UserID    int64
	Name      string
	DocUrl    string
	DocSize   int64
	CreatedAt time.Time
}
