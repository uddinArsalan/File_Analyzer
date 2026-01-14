package qdrant

import (
	"context"
	"fmt"
	"os"

	"github.com/qdrant/go-client/qdrant"
)

type QdrantClient struct {
	client         *qdrant.Client
	collectionName string
}

func NewQdrantClient(ctx context.Context, collection string) (*QdrantClient, error) {
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
	return &QdrantClient{client: qClient, collectionName: collection}, nil
}

func (qClient *QdrantClient) Close() error {
	return qClient.client.Close()
}

func (qClient *QdrantClient) CreateCollection(ctx context.Context) error {
	return qClient.client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: qClient.collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     1536,
			Distance: qdrant.Distance_Cosine,
		}),
	})
}

func (qClient *QdrantClient) EnsurePayLoadIndex(ctx context.Context, fieldName string) error {
	_, err := qClient.client.CreateFieldIndex(ctx, &qdrant.CreateFieldIndexCollection{
		CollectionName: qClient.collectionName,
		FieldName:      fieldName,
		FieldType:      qdrant.FieldType_FieldTypeKeyword.Enum(),
	})
	return err
}

func (qClient *QdrantClient) CollectionExists(ctx context.Context) (bool, error) {
	isExists, err := qClient.client.CollectionExists(ctx, qClient.collectionName)
	if err != nil {
		return false, err
	}
	return isExists, nil
}

func (qClient *QdrantClient) InsertVectorEmbeddings(ctx context.Context, points []*qdrant.PointStruct) (*qdrant.UpdateResult, error) {
	res, err := qClient.client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: qClient.collectionName,
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
	res, err := qClient.client.Query(ctx, &qdrant.QueryPoints{
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
