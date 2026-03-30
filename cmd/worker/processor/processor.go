package processor

import (
	"context"
	"file-analyzer/internals/adapters/backblaze"
	"file-analyzer/internals/adapters/cohere"
	"file-analyzer/internals/adapters/qdrant"
	"file-analyzer/internals/adapters/redis"
	"file-analyzer/internals/chunker"
	"file-analyzer/internals/domain"
	"file-analyzer/internals/parser"
	repo "file-analyzer/internals/repository"
	"file-analyzer/queue"
	"log"
	"mime"
	"strings"
)

type Processor struct {
	job    queue.Job
	llm    cohere.Embedder
	vector qdrant.VectorStore
	users  repo.UserRepository
	object backblaze.S3Store
	cache  redis.CacheStore
}

func NewProcessor(job queue.Job, llm cohere.Embedder, vector qdrant.VectorStore, users repo.UserRepository, object backblaze.S3Store, cache redis.CacheStore) *Processor {
	return &Processor{
		job,
		llm,
		vector,
		users,
		object,
		cache,
	}
}

func (p *Processor) Process(ctx context.Context, l *log.Logger) error {
	stream, err := p.object.GetObjectStream(ctx, p.job.ObjectKey)
	if err != nil {
		l.Printf("Error reading stream of file... (skipping) %v", err.Error())
		return err
	}

	// 1. Parsing
	pm := parser.NewParserManager(stream, p.job.Size)
	exts, err := mime.ExtensionsByType(p.job.Mime_Type)
	if err != nil || len(exts) == 0 {
		l.Printf("Unsupported MIME type: %s", p.job.Mime_Type)
		return err
	}
	extension := strings.TrimPrefix(exts[0], ".")
	content, err := pm.ParseFile(extension)

	if err != nil {
		l.Printf("Parsing failed: %v", err)
		return err
	}

	// l.Printf("Parsed content length: %d", len(content.Content))

	// 2. Chunking
	chunker := chunker.NewChunker(content.Content, p.job.DocID, p.job.UserID)
	chunks := chunker.Chunk()
	// l.Printf("Chunks: %+v\n", chunks)

	// 3. Embedding

	points, err := p.llm.ProcessChunks(ctx, chunks)
	if err != nil {
		l.Printf("Embedding failed: %v", err)
		return err
	}
	// l.Printf("Points: %+v\n", points)

	// 3. Adding in Vector Store each embedding wit
	_, err = p.vector.InsertVectorEmbeddings(ctx, points)
	if err != nil {
		l.Printf("Error Inserting embeddings: %v", err)
		return err
	}

	// After processing the file:
	// 1. Update the document status in the database.
	err = p.users.UpdateDocStatus(p.job.DocID, "PROCESSED")
	if err != nil {
		l.Printf("Error Updating doc status: %v", err)
		return err
	}
	// 2. Publish a "file_processed" event to Redis.
	event := domain.DocEvent{
		DocID:  p.job.DocID,
		Status: "COMPLETED",
		UserID: p.job.UserID,
	}
	err = p.cache.PublishEvent(ctx, event)
	if err != nil {
		l.Printf("Error Publishing event: %v", err)
	}

	// The API server listens for this event and notifies the frontend via SSE,
	// so the user knows the file is ready and can start asking questions.
	//
	// Optionally, send an email notification to the user once processing is complete.
	// (Especially useful for large files that take longer to process.)
	//
	// Sending an email for every file might not be ideal.
	// Need to reconsider.
	return nil
}
