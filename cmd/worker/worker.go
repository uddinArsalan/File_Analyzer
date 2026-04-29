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
	"math/rand/v2"
	"sync"
	"time"
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
	w.wg.Add(1)
	w.l.Printf("Worked ID started (MAIN WORKER) %d", w.ID)
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
					w.ProcessJobs(workerName)
				}
			}
		}
	}()
}

const MAX_RETRY_COUNT = 3
const BASE_DELAY = 1000 // 1 sec

func (w *Worker) ProcessJobs(workerName string) {
	jobs, err := w.cache.ReadJobByConsumer(w.ctx, workerName)
	if err != nil {
		w.l.Printf("Error reading jobs from stream")
		return
	}
	for _, job := range jobs {
		// process job
		w.l.Printf("Processing job %v for worker %v", job, workerName)
		processor := processor.NewProcessor(job, w.llm, w.vector, w.users, w.object, w.cache)
		err = processor.Process(w.ctx, w.l)
		if err != nil {
			w.l.Printf("Error Processing job %v err = %v", job, err)
			// But send ack still so it removes from pel
			retryCount := job.RetryCount
			if retryCount < MAX_RETRY_COUNT {
				err := w.users.UpdateDocStatus(job.DocID, "RETRYING")
				if err != nil {
					w.l.Printf("Error Updating status of doc with id %v err = %v", job.DocID, err)
					continue
				}
				base := BASE_DELAY * (1 << retryCount)
				jitter := rand.IntN(300)

				backoff := base + jitter
				retryAt := time.Now().Add(time.Duration(backoff * int(time.Second))).Unix()
				job.RetryCount++
				// will add job id here
				// and job payload in HSET and then in retry worker
				// get job from Hset and prepare job payload then
				// ZPOPMin job id and corresponding job from above and add in main stream
				// It should be in lua script so that its safe if worker
				// crash in between after removing job and before adding in
				// main queue job will be lost forever
				// and ZPopMin is atomic if multiple workers try simultaneously
				err = w.cache.AddJobToSortedSet(w.ctx, job.ID, float64(retryAt))
				if err != nil {
					w.l.Printf("Error adding job to redis sorted set %v", err)
					continue
				}

			} else {
				// Move into dead letter queue
				err := w.cache.EnqueueJobToDeadLetterQueue(w.ctx, job)
				if err != nil {
					w.l.Printf("Error adding job to dead letter queue %v", err)
					continue
				}
				// Update status
				err = w.users.UpdateDocStatus(job.DocID, "FAILED")
				if err != nil {
					w.l.Printf("Error Updating status of doc with id %v", job.DocID)
					continue
				}
			}
			// A worker will pick this job from sorted set and add it in main queue
		}
		if err := w.cache.SendAck(w.ctx, job.StreamID); err != nil {
			w.l.Printf("Error sending Acknowledgement.. %v", err.Error())
			return
		}
	}
}
