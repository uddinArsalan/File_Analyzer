package queue

type Job struct {
	ID        string
	ObjectKey string
	UserID    int64
	DocID     string
}