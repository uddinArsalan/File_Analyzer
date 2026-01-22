package services

import (
	"context"
	"file-analyzer/internals/adapters/backblaze"
	llm "file-analyzer/internals/adapters/cohere"
	"file-analyzer/internals/adapters/qdrant"
	"io"
	"mime/multipart"

	"strings"

	"github.com/google/uuid"
)

type FileService struct {
	vector   qdrant.VectorStore
	llm      llm.Embedder
	s3Client backblaze.S3Store
}

func NewFileService(vector qdrant.VectorStore,
	llm llm.Embedder, s3Client backblaze.S3Store) *FileService {
	return &FileService{
		vector:   vector,
		llm:      llm,
		s3Client: s3Client,
	}
}

func (f *FileService) UploadAndProcess(ctx context.Context, file multipart.File, userId string) (string, error) {
	// generate doc id (unique) for each document
	docId := uuid.New().String()

	const (
		maxChunks  = 96
		chunkSize  = 400
		bufferSize = 4096
	)

	var (
		buff        = make([]byte, bufferSize)
		builder     strings.Builder
		chunkBuffer []string
	)
	for {
		n, err := file.Read(buff)
		if n > 0 {
			builder.Write(buff[:n])
			if builder.Len() >= chunkSize {
				chunkBuffer = append(chunkBuffer, builder.String())
				builder.Reset()
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
	}
	if builder.Len() > 0 {
		chunkBuffer = append(chunkBuffer, builder.String())
		builder.Reset()
	}

	for i := 0; i < len(chunkBuffer); i += maxChunks {
		end := min(i+maxChunks, len(chunkBuffer))
		chunks := chunkBuffer[i:end]
		points, err := f.llm.ProcessChunks(ctx, userId, docId, chunks)
		if err != nil {
			return "", err
		}
		// after it store in db (doc id , user id ,file meta info ) maybe
		if _, err := f.vector.InsertVectorEmbeddings(ctx, points); err != nil {
			return "", err
		}
	}
	return docId, nil
}

func (f *FileService) GeneratePresignedURL(ctx context.Context, userId string, fileName string) (string, error) {
	url, err := f.s3Client.GeneratePresignedURL(ctx, userId, fileName)
	if err != nil {
		return "", nil
	}
	return url, nil
}
