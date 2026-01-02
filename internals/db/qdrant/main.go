package db

import (
	"context"
	"fmt"
	"os"

	"github.com/qdrant/go-client/qdrant"
)

type QdrantClient struct {
	Qdrant *qdrant.Client
}

func NewQdrantClient(ctx context.Context) (*QdrantClient, error) {
	qClient, err := qdrant.NewClient(&qdrant.Config{
		Host:     os.Getenv("QDRANT_HOST"),
		APIKey:   os.Getenv("QDRANT_API_KEY"),
		UseTLS:   true,
		PoolSize: 5,
	})
	if err != nil {
		return nil, fmt.Errorf("qdrant init error %w", err)
	}

	_, err = qClient.HealthCheck(ctx)
	if err != nil {
		return nil, fmt.Errorf("qdrant unreachable: %w", err)
	}
	return &QdrantClient{Qdrant: qClient}, nil
}

func (qClient *QdrantClient) Close() error {
	return qClient.Qdrant.Close()
}

func (qClient *QdrantClient) CreateCollection(ctx context.Context) error {
	return qClient.Qdrant.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: "documents",
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     1536,
			Distance: qdrant.Distance_Cosine,
		}),
	})
}

func (qClient *QdrantClient) EnsurePayLoadIndex(ctx context.Context, fieldName string) error {
	_, err := qClient.Qdrant.CreateFieldIndex(ctx, &qdrant.CreateFieldIndexCollection{
		CollectionName: "documents",
		FieldName:      fieldName,
		FieldType : qdrant.FieldType_FieldTypeKeyword.Enum(),
	})
	return err
}

func (qClient *QdrantClient) CollectionExists(ctx context.Context) (bool, error) {
	isExists, err := qClient.Qdrant.CollectionExists(ctx, "documents")
	if err != nil {
		return false, err
	}
	return isExists, nil
}

func (qClient *QdrantClient) InsertVectorEmbeddings(points []*qdrant.PointStruct) (*qdrant.UpdateResult, error) {
	res, err := qClient.Qdrant.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: "documents",
		Points:         points,
	})
	return res, err
}

func (qClient *QdrantClient) SearchEmbedInDocument(ctx context.Context, embedding []float64, docId string) ([]*qdrant.ScoredPoint, error) {
	var embed = make([]float32, len(embedding))
	for i, val := range embedding {
		embed[i] = float32(val)
	}
	// threshold := float32(0.75)
	res, err := qClient.Qdrant.Query(ctx, &qdrant.QueryPoints{
		CollectionName: "documents",
		Query:          qdrant.NewQueryDense(embed),
		Filter: &qdrant.Filter{
			Must: []*qdrant.Condition{
				qdrant.NewMatch("doc_id", docId),
			},
		},
		WithPayload: qdrant.NewWithPayload(true),
		WithVectors: qdrant.NewWithVectorsInclude(),
		// ScoreThreshold: &threshold,
	})
	return res, err
}
