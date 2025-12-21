package analytics

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaStream struct {
	writer *kafka.Writer
	topic  string
}

func NewKafkaStream(config map[string]string) (*KafkaStream, error) {
	brokers := config["brokers"]
	if brokers == "" {
		brokers = "localhost:9092"
	}

	topic := config["topic"]
	if topic == "" {
		topic = "game-analytics"
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	log.Printf("Kafka stream initialized: brokers=%s topic=%s", brokers, topic)

	return &KafkaStream{
		writer: writer,
		topic:  topic,
	}, nil
}

func (k *KafkaStream) PublishGameCompleted(gameID, winner string, duration int64) error {
	event := &GameEvent{
		Event:     "GAME_END",
		GameID:    gameID,
		Winner:    winner,
		Duration:  duration,
		Timestamp: time.Now(),
	}

	jsonBytes, err := event.Marshal()
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return err
	}

	err = k.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(gameID),
			Value: jsonBytes,
		},
	)

	if err != nil {
		log.Printf("Failed to write to Kafka: %v", err)
		return err
	}

	return nil
}

func (k *KafkaStream) Close() error {
	if k.writer != nil {
		return k.writer.Close()
	}
	return nil
}
