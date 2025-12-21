package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"4-in-a-row/analytics"
	"4-in-a-row/config"
	"4-in-a-row/db"

	"github.com/redis/go-redis/v9"
)

func main() {
	log.Println("Starting Redis Analytics Consumer...")

	// Load configuration
	cfg := config.Load()

	// Initialize database
	if err := db.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Redis configuration from config
	redisURL := cfg.RedisURL
	redisPassword := cfg.RedisPassword
	streamName := cfg.RedisStream

	groupName := "analytics-group"
	consumerName := "consumer-1"

	var opts *redis.Options

	// Try parsing as full URL first (for cloud: redis://, rediss://)
	if len(redisURL) >= 8 && (redisURL[:8] == "redis://" || (len(redisURL) >= 9 && redisURL[:9] == "rediss://")) {
		var err error
		opts, err = redis.ParseURL(redisURL)
		if err != nil {
			log.Fatal("Failed to parse Redis URL:", err)
		}
	} else {
		// Simple host:port format (for local Redis)
		opts = &redis.Options{
			Addr:     redisURL,
			Password: redisPassword,
			DB:       0,
		}
	}

	// Override password if provided separately (only for local)
	if redisPassword != "" && opts.Password == "" {
		opts.Password = redisPassword
	}

	client := redis.NewClient(opts)
	defer client.Close()

	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	log.Println("Connected to Redis successfully")

	// Create consumer group (ignore error if exists)
	client.XGroupCreateMkStream(ctx, streamName, groupName, "0")
	log.Printf("Redis consumer started (stream: %s, group: %s)", streamName, groupName)

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Consumer loop
	go func() {
		for {
			// Read from stream with blocking
			streams, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    groupName,
				Consumer: consumerName,
				Streams:  []string{streamName, ">"},
				Count:    10,
				Block:    5 * time.Second,
			}).Result()

			if err != nil && err != redis.Nil {
				log.Println("Error reading from stream:", err)
				time.Sleep(time.Second)
				continue
			}

			// Process messages
			for _, stream := range streams {
				for _, message := range stream.Messages {
					processMessage(ctx, client, streamName, groupName, message)
				}
			}
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down Redis consumer...")
}

func processMessage(ctx context.Context, client *redis.Client, streamName, groupName string, msg redis.XMessage) {
	// Parse event from fields
	event := analytics.GameEvent{
		Event:  getStringField(msg.Values, "event"),
		GameID: getStringField(msg.Values, "gameId"),
		Winner: getStringField(msg.Values, "winner"),
	}

	// Parse duration
	if durationStr, ok := msg.Values["duration"].(string); ok {
		var duration int64
		json.Unmarshal([]byte(durationStr), &duration)
		event.Duration = duration
	}

	// Parse timestamp
	if timestampStr, ok := msg.Values["timestamp"].(string); ok {
		if parsedTime, err := time.Parse(time.RFC3339, timestampStr); err == nil {
			event.Timestamp = parsedTime
		}
	}

	// Process the event
	processGameEvent(event)

	// Acknowledge message
	client.XAck(ctx, streamName, groupName, msg.ID)
}

func getStringField(values map[string]interface{}, key string) string {
	if val, ok := values[key].(string); ok {
		return val
	}
	return ""
}

func processGameEvent(event analytics.GameEvent) {
	log.Printf("Processing event: %s for game %s", event.Event, event.GameID)

	if event.Event == "GAME_END" {
		// Update game metrics (hourly aggregation)
		if err := db.UpdateGameMetrics(event.Duration, event.Timestamp); err != nil {
			log.Printf("Failed to update game metrics: %v", err)
		}

		// Update player statistics
		isDraw := event.Winner == "draw"

		// Get player names from the game result
		var gameResult db.GameResult
		if err := db.DB.Where("game_id = ?", event.GameID).First(&gameResult).Error; err != nil {
			log.Printf("Failed to fetch game result: %v", err)
			return
		}

		// Update player 1 stats
		player1Won := !isDraw && event.Winner == gameResult.Player1.Username
		if err := db.UpdatePlayerStats(gameResult.Player1.Username, player1Won, isDraw, event.Duration); err != nil {
			log.Printf("Failed to update player 1 stats: %v", err)
		}

		// Update player 2 stats (skip if bot)
		if gameResult.Player2.Type != "bot" {
			player2Won := !isDraw && event.Winner == gameResult.Player2.Username
			if err := db.UpdatePlayerStats(gameResult.Player2.Username, player2Won, isDraw, event.Duration); err != nil {
				log.Printf("Failed to update player 2 stats: %v", err)
			}
		}
	}
}
