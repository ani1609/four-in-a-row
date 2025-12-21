package analytics

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStream struct {
	client     *redis.Client
	streamName string
	ctx        context.Context
}

func NewRedisStream(config map[string]string) (*RedisStream, error) {
	redisURL := config["url"]
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	password := config["password"]
	streamName := config["stream"]
	if streamName == "" {
		streamName = "game-events"
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		opt = &redis.Options{
			Addr:     redisURL,
			Password: password,
			DB:       0,
		}
	}

	if password != "" {
		opt.Password = password
	}

	client := redis.NewClient(opt)
	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	log.Printf("Redis stream initialized: addr=%s stream=%s", opt.Addr, streamName)

	return &RedisStream{
		client:     client,
		streamName: streamName,
		ctx:        ctx,
	}, nil
}

func (r *RedisStream) PublishGameCompleted(gameID, winner string, duration int64) error {
	event := &GameEvent{
		Event:     "GAME_END",
		GameID:    gameID,
		Winner:    winner,
		Duration:  duration,
		Timestamp: time.Now(),
	}

	// Convert event to Redis Stream fields
	fields := map[string]interface{}{
		"event":     event.Event,
		"gameId":    event.GameID,
		"winner":    event.Winner,
		"duration":  event.Duration,
		"timestamp": event.Timestamp.Format(time.RFC3339),
	}

	id, err := r.client.XAdd(r.ctx, &redis.XAddArgs{
		Stream: r.streamName,
		Values: fields,
	}).Result()

	if err != nil {
		log.Printf("Failed to write to Redis stream: %v", err)
		return err
	}

	log.Printf("Published event to Redis stream: id=%s gameId=%s", id, gameID)
	return nil
}

func (r *RedisStream) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}
