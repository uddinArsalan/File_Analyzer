package worker

import (
	"context"
	"file-analyzer/cmd/worker/processor"
	"file-analyzer/internals/adapters/backblaze"
	"file-analyzer/internals/adapters/cohere"
	"file-analyzer/internals/adapters/qdrant"
	repo "file-analyzer/internals/repository"
	"file-analyzer/queue"
	"log"
	"time"
)

type Worker struct {
	ID      int
	JobChan chan queue.Job
	ctx     context.Context
	l       *log.Logger
	llm     cohere.Embedder
	vector  qdrant.VectorStore
	users   repo.UserRepository
	object  backblaze.S3Store
}

func (w *Worker) Start() {
	processor := processor.NewProcessor(w.llm, w.vector, w.users, w.object)
	w.l.Printf("Worked ID started %d", w.ID)
	go func() {
		for {
			select {
			case <-w.ctx.Done():
				{
					w.l.Printf("Worker #%d stopping...\n", w.ID)
					return
				}
			case job, ok := <-w.JobChan:
				{
					if !ok {
						w.l.Printf("Job channel closed, worker %d exiting", w.ID)
						return
					}

					jobCtx, cancel := context.WithTimeout(w.ctx, 10*time.Minute)
					err := processor.Process(jobCtx, job)
					cancel()
					if err != nil {
						w.l.Printf("Worker %d failed job %s: %v", w.ID, job.ID, err)

					} else {
						w.l.Printf("Worker %d completed job %s", w.ID, job.ID)
					}
				}
			}
		}
	}()
}
