package analytics

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

var Writer *kafka.Writer

func InitKafka(brokers []string, topic string) {
	Writer = &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	log.Println("Kafka Writer initialized")
}

type GameEvent struct {
	Event     string    `json:"event"`
	GameID    string    `json:"gameId"`
	Winner    string    `json:"winner"`
	Duration  int64     `json:"duration"`
	Timestamp time.Time `json:"timestamp"`
}

func EmitGameEnd(gameID, winner string, duration int64) {
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
