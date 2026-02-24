package main

import (
	"context"
	"file-analyzer/internals/adapters/backblaze"
	"file-analyzer/internals/adapters/cohere"
	"file-analyzer/internals/adapters/qdrant"
	"file-analyzer/internals/adapters/redis"
	repo "file-analyzer/internals/repository"
	"log"
	"sync"
)

type Dispatcher struct {
	WorkerCount int
	workers     []*Worker
	ctx         context.Context
	cancel      context.CancelFunc
	wg          *sync.WaitGroup
}

func NewDispatcher(parent context.Context, workerCount, queueSize int) *Dispatcher {
	ctx, cancel := context.WithCancel(parent)

	return &Dispatcher{
		WorkerCount: workerCount,
		ctx:         ctx,
		cancel:      cancel,
		wg:          &sync.WaitGroup{},
	}
}

func (d *Dispatcher) Start(l *log.Logger, llm cohere.Embedder, vector qdrant.VectorStore, users repo.UserRepository, object backblaze.S3Store, cache redis.CacheStore) {
	for i := 1; i <= d.WorkerCount; i++ {
		d.wg.Add(1)
		worker := &Worker{
			ID:     i,
			ctx:    d.ctx,
			l:      l,
			llm:    llm,
			vector: vector,
			users:  users,
			object: object,
			wg:     d.wg,
			cache:  cache,
		}
		d.workers = append(d.workers, worker)
		worker.Start()
	}
}

func (d *Dispatcher) Stop() {
	d.cancel()  // stop workers
	d.wg.Wait() // wait for workers to finish
}

// func (d *Dispatcher) StartRedisListener(ctx context.Context, l *log.Logger, cache redis.CacheStore) {
// 	lastID := "$"

// 	go func() {
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				{
// 					l.Println("Redis listener stopping...")
// 					return
// 				}
// 			default:
// 				{
// 					streams, err := cache.DequeueJob(ctx, lastID)

// 					if err != nil {
// 						l.Println("Redis read error:", err)
// 						time.Sleep(2 * time.Second)
// 						continue
// 					}

// 					for _, stream := range streams {
// 						for _, msg := range stream.Messages {
// 							lastID = msg.ID
// 							userIDStr := msg.Values["user_id"].(string)
// 							userID, err := strconv.ParseInt(userIDStr, 10, 64)
// 							if err != nil {
// 								l.Println("invalid user_id:", err)
// 								continue
// 							}
// 							job := queue.Job{
// 								ID:        msg.Values["id"].(string),
// 								UserID:    userID,
// 								ObjectKey: msg.Values["object_key"].(string),
// 								DocID:     msg.Values["doc_id"].(string),
// 							}
// 							d.Submit(job)
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}()
// }
