package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	// "file-analyzer/cmd/worker/processor"
	"file-analyzer/internals/adapters/backblaze"
	"file-analyzer/internals/adapters/cohere"
	"file-analyzer/internals/adapters/qdrant"
	"file-analyzer/internals/adapters/redis"
	repo "file-analyzer/internals/repository"
	"file-analyzer/queue"
	"log"
	"sync"
	// "time"
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
	// processor := processor.NewProcessor(w.llm, w.vector, w.users, w.object)
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
					w.l.Printf("JOB PICKED BY %v", workerName)
					streams, err := w.cache.ReadJobByConsumer(w.ctx, workerName)

					if err != nil {
						w.l.Println("Redis read error:", err)
						time.Sleep(2 * time.Second)
						continue
					}

					for _, stream := range streams {
						for _, msg := range stream.Messages {
							userIDStr := msg.Values["user_id"].(string)
							userID, err := strconv.ParseInt(userIDStr, 10, 64)
							if err != nil {
								w.l.Println("invalid user_id:", err)
								continue
							}
							job := queue.Job{
								ID:        msg.Values["id"].(string),
								UserID:    userID,
								ObjectKey: msg.Values["object_key"].(string),
								DocID:     msg.Values["doc_id"].(string),
							}
							// process job
							w.l.Printf("JOB %v", job)
							// if success
							// send acknowledgement use XACK
							if err := w.cache.SendAck(w.ctx, msg.ID); err != nil {
								w.l.Printf("Error sending Acknowledgement.. %v", err.Error())
							}
						}
					}
				}
			}
		}
	}()
}
