package queue

type Job struct {
	ID        string
	ObjectKey string
	UserID    int64
	DocID     string
	Mime_Type string
	Size      int64
}
