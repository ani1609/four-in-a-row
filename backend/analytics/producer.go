package analytics

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// Legacy Kafka writer (kept for backward compatibility)
var Writer *kafka.Writer

// InitKafka initializes legacy Kafka writer (deprecated - use InitEventStream)
func InitKafka(brokers []string, topic string) {
	Writer = &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	log.Println("Kafka Writer initialized (legacy)")
}

// EmitGameEnd emits a game end event (uses new event stream abstraction)
func EmitGameEnd(gameID, winner string, duration int64) {
	// Try new abstraction first
	if err := PublishGameCompleted(gameID, winner, duration); err == nil {
		return
	}

	// Fallback to legacy Kafka writer
	if Writer == nil {
		return
	}

	event := GameEvent{
		Event:     "GAME_END",
		GameID:    gameID,
		Winner:    winner,
		Duration:  duration,
		Timestamp: time.Now(),
	}

	jsonBytes, _ := json.Marshal(event)

	err := Writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(gameID),
			Value: jsonBytes,
		},
	)
	if err != nil {
		log.Println("Failed to write to Kafka:", err)
	}
}
