package processor

import (
	"bufio"
	"context"
	"file-analyzer/internals/adapters/backblaze"
	"file-analyzer/internals/adapters/cohere"
	"file-analyzer/internals/adapters/qdrant"
	repo "file-analyzer/internals/repository"
	"file-analyzer/queue"
	"io"
	"strings"
)

type Processor struct {
	llm    cohere.Embedder
	vector qdrant.VectorStore
	users  repo.UserRepository
	object backblaze.S3Store
}

func NewProcessor(llm cohere.Embedder, vector qdrant.VectorStore, users repo.UserRepository, object backblaze.S3Store) *Processor {
	return &Processor{
		llm,
		vector,
		users,
		object,
	}
}

func (p *Processor) Process(ctx context.Context, job queue.Job) error {
	err := p.users.UpdateDocStatus(job.DocID, "PROCESSING")
	if err != nil {
		return nil
	}
	body, err := p.object.GetObjectStream(ctx, job.ObjectKey)
	if err != nil {
		return err
	}
	return p.UploadAndProcess(ctx, job, body)
}

func (p *Processor) UploadAndProcess(ctx context.Context, job queue.Job, body io.ReadCloser) error {

	reader := bufio.NewReader(body)

	const (
		maxChunks  = 96
		chunkSize  = 400
		bufferSize = 4096
	)

	var (
		buff        = make([]byte, bufferSize)
		builder     strings.Builder
		chunkBuffer []string
	)

	defer body.Close()

	for {
		n, err := reader.Read(buff)
		if n > 0 {
			builder.Write(buff[:n])
			if builder.Len() >= chunkSize {
				chunkBuffer = append(chunkBuffer, builder.String())
				builder.Reset()
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	if builder.Len() > 0 {
		chunkBuffer = append(chunkBuffer, builder.String())
		builder.Reset()
	}

	for i := 0; i < len(chunkBuffer); i += maxChunks {
		end := min(i+maxChunks, len(chunkBuffer))
		chunks := chunkBuffer[i:end]
		points, err := p.llm.ProcessChunks(ctx, job.UserID, job.DocID, chunks)
		if err != nil {
			return err
		}
		if _, err := p.vector.InsertVectorEmbeddings(ctx, points); err != nil {
			return err
		}
	}
	return nil
}
