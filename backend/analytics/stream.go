package analytics

import (
	"encoding/json"
	"time"
)

type EventStream interface {
	PublishGameCompleted(gameID, winner string, duration int64) error
	Close() error
}

type GameEvent struct {
	Event     string    `json:"event"`
	GameID    string    `json:"gameId"`
	Winner    string    `json:"winner"`
	Duration  int64     `json:"duration"`
	Timestamp time.Time `json:"timestamp"`
}

func (e *GameEvent) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func UnmarshalGameEvent(data []byte) (*GameEvent, error) {
	var event GameEvent
	err := json.Unmarshal(data, &event)
	return &event, err
}

var globalStream EventStream

func InitEventStream(streamType string, config map[string]string) error {
	var stream EventStream
	var err error

	switch streamType {
	case "redis":
		stream, err = NewRedisStream(config)
	case "kafka":
		stream, err = NewKafkaStream(config)
	default:
		stream, err = NewKafkaStream(config)
	}

	if err != nil {
		return err
	}

	globalStream = stream
	return nil
}

func PublishGameCompleted(gameID, winner string, duration int64) error {
	if globalStream == nil {
		return nil
	}
	return globalStream.PublishGameCompleted(gameID, winner, duration)
}

func Close() error {
	if globalStream != nil {
		return globalStream.Close()
	}
	return nil
}
