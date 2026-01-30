package backblaze

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	ps *s3.PresignClient
	sc *s3.Client
}

func NewS3Client(ctx context.Context) (*S3Client, error) {
	endpoint := os.Getenv("ENDPOINT")
	region := os.Getenv("REGION")

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(
			aws.NewCredentialsCache(
				aws.CredentialsProviderFunc(
					func(ctx context.Context) (aws.Credentials, error) {
						return aws.Credentials{
							AccessKeyID:     os.Getenv("KEY_ID"),
							SecretAccessKey: os.Getenv("APP_KEY"),
						}, nil
					},
				),
			),
		),
		config.WithDefaultRegion(region),
	)

	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})
	presignClient := s3.NewPresignClient(client)
	return &S3Client{ps: presignClient, sc: client}, nil
}

func (s3C *S3Client) GeneratePresignedURL(ctx context.Context, objectKey string) (string, error) {
	bucketName := os.Getenv("BUCKET_NAME")
	req, err := s3C.ps.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, s3.WithPresignExpires(5*time.Minute))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (s3C *S3Client) GetObjectStream(ctx context.Context, key string) (io.ReadCloser, error) {
	bucketName := os.Getenv("BUCKET_NAME")

	res, err := s3C.sc.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}
	return res.Body, nil
}
