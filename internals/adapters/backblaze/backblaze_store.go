package backblaze

import "context"

type S3Store interface {
	GeneratePresignedURL(ctx context.Context, userId string, fileName string) (string, error)
}
