package main

import (
	"encoding/json"
	"file-analyzer/queue"
	"time"
)

func (w *Worker) StartRetryingJobs() {
	go func() {
		defer w.wg.Done()
		for {
			ticker := time.NewTicker(30 * time.Second)
			select {
			case <-w.ctx.Done():
				{
					w.l.Printf("Worker #%d stopping...\n", w.ID)
					return
				}
			case <-ticker.C:
				{
					res, err := w.cache.GetJobsReadyForRetry(w.ctx)
					if err != nil {
						w.l.Printf("Error getting jobs to retry err = %v\n", err)
						return
					}
					for _, item := range res {
						var job queue.Job
						err := json.Unmarshal([]byte(item), &job)
						if err != nil {
							w.l.Printf("Error unmarshal job , err = %v", err)
							continue
						}
						err = w.cache.EnqueueJob(w.ctx, job)
						if err != nil {
							w.l.Printf("Error adding job to stream %v\n", err)
							continue
						}
					}

				}
			}
		}
	}()
}
