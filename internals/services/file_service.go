package services

import (
	"context"
	"file-analyzer/internals/adapters/backblaze"
	"file-analyzer/internals/adapters/redis"
	"file-analyzer/internals/domain"
	"file-analyzer/internals/handlers/dto"
	repo "file-analyzer/internals/repository"
	"file-analyzer/internals/sse"
	"file-analyzer/queue"
	"net/http"

	// "file-analyzer/queue"
	"fmt"

	"github.com/google/uuid"
	// "github.com/google/uuid"
)

type FileService struct {
	s3Client backblaze.S3Store
	users    repo.UserRepository
	cache    redis.CacheStore
	sse      *sse.SSEManager
}

func NewFileService(s3Client backblaze.S3Store, users repo.UserRepository, cache redis.CacheStore, sse *sse.SSEManager) *FileService {
	return &FileService{
		s3Client,
		users,
		cache,
		sse,
	}
}

// Consider idempotency
func (f *FileService) CheckExistence(ctx context.Context, userID int64, docID string) error {
	doc, err := f.users.DocumentExistsForUser(userID, docID)
	if err != nil {
		return ErrDocumentNotFound
	}
	objectKey := fmt.Sprintf("documents/%v/%v", userID, docID)
	// head request to object storage to check if file is uploaded
	isExists, err := f.s3Client.HeadObject(ctx, objectKey)
	if err != nil {
		return err
	}
	if !isExists {
		return ErrDocumentNotFound
	}
	job := &queue.Job{
		ID:        uuid.New().String(),
		ObjectKey: objectKey,
		UserID:    userID,
		DocID:     docID,
		Mime_Type: doc.Mime_Type,
		Size:      doc.DocSize,
	}
	if err := f.cache.EnqueueJob(ctx, job); err != nil {
		return err
	}
	return f.users.UpdateDocStatus(docID, "PROCESSING")
}

func (f *FileService) GeneratePresignedURL(ctx context.Context, userID int64, docID string, doc dto.DocRequest) (string, error) {
	objectKey := fmt.Sprintf("documents/%v/%v", userID, docID)
	docObj := domain.Document{
		DocID:     docID,
		UserID:    userID,
		Name:      doc.FileName,
		ObjectKey: objectKey,
		Status:    "PENDING",
		Mime_Type: doc.MiMeType,
		DocSize:   int64(doc.FileSize),
	}
	err := f.users.InsertDoc(docID, docObj)
	if err != nil {
		return "", err
	}
	url, err := f.s3Client.GeneratePresignedURL(ctx, objectKey)
	if err != nil {
		return "", nil
	}
	return url, nil
}

func (f *FileService) Stream(ctx context.Context, userID int64, flusher http.Flusher, w http.ResponseWriter) {
	ch := f.sse.AddClient(userID)

	defer f.sse.RemoveClient(userID)

	for {
		select {
		case event := <-ch:
			{
				fmt.Fprintf(w, "data: %v\n\n", event)
				flusher.Flush()
			}
		case <-ctx.Done():
			{
				return
			}
		}
	}
}
