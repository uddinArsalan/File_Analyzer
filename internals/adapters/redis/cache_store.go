package redis

import (
	"context"
	"file-analyzer/internals/domain"
	"file-analyzer/internals/subscriber"
	"file-analyzer/queue"
)

type CacheStore interface {
	EnqueueJob(ctx context.Context, job queue.Job) error
	ReadJobByConsumer(ctx context.Context, consumer string) ([]queue.Job, error)
	SendAck(ctx context.Context, id string) error
	CreateAndCheckStream(parent context.Context) error
	PublishEvent(ctx context.Context, message domain.DocEvent) error
	SubscribeAndListen(ctx context.Context, subscribers []subscriber.Subscriber) error
	GetPendingJobs(ctx context.Context) ([]XPending, error)
	ClaimPendingJobs(ctx context.Context, consumerName string, messageIds []string) ([]queue.Job, error)
	AddJobToSortedSet(ctx context.Context, job string, timestamp float64) error
	EnqueueJobToDeadLetterQueue(ctx context.Context, job queue.Job) error
	GetJobIDsReadyForRetry(ctx context.Context) ([]string, error)
}
