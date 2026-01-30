package services

import (
	"context"
	"file-analyzer/internals/adapters/backblaze"
	"file-analyzer/internals/domain"
	"file-analyzer/internals/handlers/dto"
	repo "file-analyzer/internals/repository"
	"fmt"
)

type FileService struct {
	s3Client backblaze.S3Store
	users    repo.UserRepository
}

func NewFileService(s3Client backblaze.S3Store, users repo.UserRepository) *FileService {
	return &FileService{
		s3Client: s3Client,
		users:    users,
	}
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
