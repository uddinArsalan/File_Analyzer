package main

import (
	"file-analyzer/cmd/worker/processor"
	"fmt"
	"time"
)

func (w *Worker) StartRecoveryWorker() {
	w.wg.Add(1)
	w.l.Printf("Worked ID started (RECOVERY WORKER) %d", w.ID)
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer w.wg.Done()
		for {
			select {
			case <-w.ctx.Done():
				{
					w.l.Printf("Shutting down failure recovery... %v",w.ID)
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
					if len(pendingEntryList) == 0 {
						w.l.Printf("No jobs to recover (empty pel)")
						return
					}
					pending := make([]string, len(pendingEntryList))
					for i, entry := range pendingEntryList {
						pending[i] = entry.ID
					}
					w.l.Printf("Pending jobs %v",pending)
					jobs, err := w.cache.ClaimPendingJobs(w.ctx, workerName, pending)
					if err != nil {
						w.l.Printf("Error Claiming Jobs %v", err)
					}
					for _, job := range jobs {
						w.l.Printf("Processing job %+v for worker %s", job, workerName)
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
