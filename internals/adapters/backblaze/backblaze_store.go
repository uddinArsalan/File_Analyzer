package backblaze

import (
	"context"
	"io"
)

type S3Store interface {
	GeneratePresignedURL(ctx context.Context, objectKey string) (string, error)
	GetObjectStream(ctx context.Context, key string) (io.ReadCloser, error)
	HeadObject(ctx context.Context, key string) (bool, error)
}
