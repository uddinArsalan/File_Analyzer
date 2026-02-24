package services

import (
	"context"
	"file-analyzer/internals/adapters/redis"
	"file-analyzer/queue"
	"fmt"

	"github.com/google/uuid"
)

type TempService struct {
	cache redis.CacheStore
}

func NewTempService(cache redis.CacheStore) *TempService {
	return &TempService{
		cache,
	}
}

func (t *TempService) AddTemporaryJobs(ctx context.Context, i int64) {
	for j := range i {
		t.cache.EnqueueJob(ctx, &queue.Job{
			ID:        uuid.New().String(),
			UserID:    i,
			ObjectKey: fmt.Sprintf("key:%v", j),
			DocID:     uuid.New().String(),
		})
	}
}
