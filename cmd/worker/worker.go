package main

import (
	"context"
	"file-analyzer/cmd/worker/processor"
	"file-analyzer/internals/adapters/backblaze"
	"file-analyzer/internals/adapters/cohere"
	"file-analyzer/internals/adapters/qdrant"
	"file-analyzer/internals/adapters/redis"
	repo "file-analyzer/internals/repository"
	"fmt"
	"log"
	"sync"
)

type Worker struct {
	ID     int
	ctx    context.Context
	l      *log.Logger
	llm    cohere.Embedder
	vector qdrant.VectorStore
	users  repo.UserRepository
	object backblaze.S3Store
	cache  redis.CacheStore
	wg     *sync.WaitGroup
}

func (w *Worker) Start() {
	w.l.Printf("Worked ID started %d", w.ID)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case <-w.ctx.Done():
				{
					w.l.Printf("Worker #%d stopping...\n", w.ID)
					return
				}
			default:
				{
					workerName := fmt.Sprintf("Worker #%d", w.ID)
					w.l.Printf("Job picked by %s", workerName)

					jobs, err := w.cache.ReadJobByConsumer(w.ctx, workerName)
					if err != nil {
						w.l.Printf("Error reading jobs from stream")
						break
					}
					for _, job := range jobs {
						// process job
						w.l.Printf("Processing job %v for worker %v", job, workerName)
						processor := processor.NewProcessor(job, w.llm, w.vector, w.users, w.object, w.cache)
						err = processor.Process(w.ctx, w.l)
						if err != nil {
							continue
						}
						// if success
						// send acknowledgement use XACK
						if err := w.cache.SendAck(w.ctx, job.ID); err != nil {
							w.l.Printf("Error sending Acknowledgement.. %v", err.Error())
						}
					}

				}
			}
		}
	}()
}
