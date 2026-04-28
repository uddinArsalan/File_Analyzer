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
	l           *log.Logger
	llm         cohere.Embedder
	vector      qdrant.VectorStore
	users       repo.UserRepository
	object      backblaze.S3Store
	cache       redis.CacheStore
}

func NewDispatcher(parent context.Context, workerCount, queueSize int, l *log.Logger, llm cohere.Embedder, vector qdrant.VectorStore, users repo.UserRepository, object backblaze.S3Store, cache redis.CacheStore) *Dispatcher {
	ctx, cancel := context.WithCancel(parent)

	return &Dispatcher{
		WorkerCount: workerCount,
		ctx:         ctx,
		cancel:      cancel,
		wg:          &sync.WaitGroup{},
		l:           l,
		llm:         llm,
		vector:      vector,
		users:       users,
		object:      object,
		cache:       cache,
	}
}

func (d *Dispatcher) Start() {
	for i := 1; i <= d.WorkerCount; i++ {
		d.wg.Add(1)
		worker := &Worker{
			ID:     i,
			ctx:    d.ctx,
			l:      d.l,
			llm:    d.llm,
			vector: d.vector,
			users:  d.users,
			object: d.object,
			wg:     d.wg,
			cache:  d.cache,
		}
		d.workers = append(d.workers, worker)
		worker.Start()
		worker.StartRecoveryWorker()
		worker.StartRetryingJobs()
	}
}

func (d *Dispatcher) Stop() {
	d.cancel()  // stop workers
	d.wg.Wait() // wait for workers to finish
}
