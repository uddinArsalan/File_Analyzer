package backblaze

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	ps *s3.PresignClient
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
	return &S3Client{ps: presignClient}, nil
}

func (s3C *S3Client) GeneratePresignedURL(ctx context.Context, userId string, fileName string) (string, error) {
	key := fmt.Sprintf("documents/%v/%v", userId, fileName)
	bucketName := os.Getenv("BUCKET_NAME")
	req, err := s3C.ps.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:  aws.String(bucketName),
		Key:     aws.String(key),
		Expires: aws.Time(time.Now().Add(5 * time.Minute)),
	}, s3.WithPresignExpires(5*time.Minute))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}
