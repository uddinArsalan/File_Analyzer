package main

import (
	"context"
	"file-analyzer/cmd/worker/processor"
	"file-analyzer/internals/adapters/backblaze"
	"file-analyzer/internals/adapters/cohere"
	"file-analyzer/internals/adapters/qdrant"
	"file-analyzer/internals/adapters/redis"
	repo "file-analyzer/internals/repository"
	"file-analyzer/queue"
	"fmt"
	"log"
	"strconv"
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
					streams, err := w.cache.ReadJobByConsumer(w.ctx, workerName)

					if err != nil {
						w.l.Println("Redis read error:", err)
						time.Sleep(2 * time.Second)
						continue
					}

					if len(streams) == 0 {
						time.Sleep(500 * time.Millisecond)
						continue
					}

					w.l.Printf("Job picked by %s", workerName)

					for _, stream := range streams {
						for _, msg := range stream.Messages {
							userIDStr, ok := msg.Values["user_id"].(string)
							if !ok {
								w.l.Printf("Worker #%d: missing or invalid user_id in msg %s", w.ID, msg.ID)
								w.cache.SendAck(w.ctx, msg.ID)
								continue
							}
							userID, err := strconv.ParseInt(userIDStr, 10, 64)
							if err != nil {
								w.l.Printf("Worker #%d: invalid user_id %q: %v", w.ID, userIDStr, err)
								w.cache.SendAck(w.ctx, msg.ID)
								continue
							}
							job := queue.Job{
								ID:        msg.Values["id"].(string),
								UserID:    userID,
								ObjectKey: msg.Values["object_key"].(string),
								DocID:     msg.Values["doc_id"].(string),
								Mime_Type: msg.Values["mime_type"].(string),
							}
							// process job
							w.l.Printf("Processing job %v for worker %v", job, workerName)
							processor := processor.NewProcessor(job, w.llm, w.vector, w.users, w.object,w.cache)
							err = processor.Process(w.ctx, w.l)
							if err != nil {
								continue
							}
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
