package queue

import (
	"context"
	"file-analyzer/internals/adapters/backblaze"
	"file-analyzer/internals/adapters/cohere"
	"file-analyzer/internals/adapters/qdrant"
	repo "file-analyzer/internals/repository"
	"log"
	"sync"
)

type Dispatcher struct {
	WorkerCount int
	JobQueue    chan Job
	workers     []*Worker
	ctx         context.Context
	cancel      context.CancelFunc
	wg          *sync.WaitGroup
}

func NewDispatcher(workerCount, queueSize int) *Dispatcher {
	ctx, cancel := context.WithCancel(context.Background())

	return &Dispatcher{
		WorkerCount: workerCount,
		JobQueue:    make(chan Job, queueSize),
		ctx:         ctx,
		cancel:      cancel,
		wg:          &sync.WaitGroup{},
	}
}

func (d *Dispatcher) Start(l *log.Logger, llm cohere.Embedder, vector qdrant.VectorStore, users repo.UserRepository, object backblaze.S3Store) {
	for i := 1; i <= d.WorkerCount; i++ {
		worker := &Worker{
			ID:      i,
			JobChan: d.JobQueue,
			ctx:     d.ctx,
			l:       l,
			llm:     llm,
			vector:  vector,
			users:   users,
			object:  object,
		}
		d.workers = append(d.workers, worker)
		worker.Start()
	}
}

func (d *Dispatcher) Submit(job Job) {
	d.wg.Add(1)
	go func() {
		d.JobQueue <- job
		d.wg.Done()
	}()
}
