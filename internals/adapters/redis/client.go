package redis

import (
	"context"
	"file-analyzer/internals/domain"
	"file-analyzer/queue"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	rdb           *redis.Client
	streamName    string
	consumerGroup string
	channelName   string
}

func NewRedisClient(ctx context.Context) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:        os.Getenv("REDIS_ADDR"),
		Password:    os.Getenv("REDIS_PSSWRD"),
		Username:    "default",
		DB:          0,
		ReadTimeout: -1,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Error Ping Connection Redis %v", err.Error())
	}
	return &RedisClient{
		rdb:           rdb,
		streamName:    os.Getenv("REDIS_STREAM"),
		consumerGroup: os.Getenv("REDIS_CONSUMER_GROUP"),
		channelName:   os.Getenv("REDIS_EVENT_CHANNEL"),
	}, nil
}

func (redisClient *RedisClient) CloseRedisClient() error {
	return redisClient.rdb.Close()
}

func (redisClient *RedisClient) EnqueueJob(ctx context.Context, job *queue.Job) error {
	values := map[string]interface{}{
		"id":         job.ID,
		"object_key": job.ObjectKey,
		"user_id":    job.UserID,
		"doc_id":     job.DocID,
		"mime_type":  job.Mime_Type,
	}
	res, err := redisClient.rdb.XAdd(ctx, &redis.XAddArgs{
		ID:     "*",
		Stream: redisClient.streamName,
		Values: values,
	}).Result()

	if err != nil {
		fmt.Printf("Error %v", err)
		return err
	}

	fmt.Printf("Document job enqueued %v \n", res)
	return nil
}

func (redisClient *RedisClient) ReadJobByConsumer(ctx context.Context, consumer string) ([]redis.XStream, error) {
	stream, err := redisClient.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    redisClient.consumerGroup,
		Consumer: consumer,
		Streams:  []string{redisClient.streamName, ">"},
		Count:    10,
		Block:    2 * time.Second,
	}).Result()

	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return stream, nil
}

func (redisClient *RedisClient) CreateAndCheckStream(parent context.Context) error {
	ctx, cancel := context.WithTimeout(parent, 20*time.Second)
	defer cancel()

	return redisClient.rdb.XGroupCreateMkStream(ctx, redisClient.streamName, redisClient.consumerGroup, "0").Err()
}

func (redisClient *RedisClient) SendAck(ctx context.Context, id string) error {
	_, err := redisClient.rdb.XAck(ctx, redisClient.streamName, redisClient.consumerGroup, id).Result()
	return err
}

func (redisClient *RedisClient) PublishEvent(ctx context.Context, message domain.DocEvent) {
	redisClient.rdb.Publish(ctx, redisClient.channelName, message)
}

func (redisClient *RedisClient) Subscribe(ctx context.Context) *redis.PubSub {
	return redisClient.rdb.Subscribe(ctx, redisClient.channelName)
}
