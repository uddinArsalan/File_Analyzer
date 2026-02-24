package redis

import (
	"context"
	"file-analyzer/queue"

	"github.com/redis/go-redis/v9"
)

type CacheStore interface {
	EnqueueJob(ctx context.Context, job *queue.Job) error
	ReadJobByConsumer(ctx context.Context, consumer string) ([]redis.XStream, error)
	SendAck(ctx context.Context, id string) error
	CreateAndCheckStream(parent context.Context) error
}
