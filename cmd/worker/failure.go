package main

import (
	"file-analyzer/cmd/worker/processor"
	"fmt"
	"time"
)

const MAX_COUNT int64 = 3
const BASE_DELAY = 100

func (w *Worker) StartRecoveryWorker() {
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer w.wg.Done()
		for {
			select {
			case <-w.ctx.Done():
				{
					w.l.Print("Shutting down failure recovery...")
					return
				}
			case <-ticker.C:
				{
					workerName := fmt.Sprintf("Worker #%d", w.ID)
					pendingEntryList, err := w.cache.GetPendingJobs(w.ctx)
					if err != nil {
						w.l.Printf("Error reading pel %v", err.Error())
						return
					}
					pending := make([]string, len(pendingEntryList))
					for _, entry := range pendingEntryList {
						pending = append(pending, entry.ID)
					}
					jobs, err := w.cache.ClaimPendingJobs(w.ctx, workerName, pending)
					if err != nil {
						w.l.Printf("Error Claiming Jobs %v", err)
					}
					for _, job := range jobs {
						w.l.Printf("Processing job %v for worker %v", job, workerName)
						processor := processor.NewProcessor(job, w.llm, w.vector, w.users, w.object, w.cache)
						err = processor.Process(w.ctx, w.l)
						if err := w.cache.SendAck(w.ctx, job.ID); err != nil {
							w.l.Printf("Error sending Acknowledgement.. %v", err.Error())
							return
						}
					}
				}
			}
		}
	}()
}
