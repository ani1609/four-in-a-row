package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"4-in-a-row/config"
	"4-in-a-row/db"

	"github.com/segmentio/kafka-go"
)

type GameEvent struct {
	Event     string    `json:"event"`
	GameID    string    `json:"gameId"`
	Winner    string    `json:"winner"`
	Duration  int64     `json:"duration"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	log.Println("Starting Analytics Consumer...")

	// Load configuration
	cfg := config.Load()

	// Initialize Database
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Analytics consumer connected to database")

	// Initialize Kafka Reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.KafkaBrokers},
		Topic:    cfg.KafkaTopic,
		GroupID:  "analytics-group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	defer r.Close()

	log.Println("Listening for game events on topic: game-analytics")

	// Process messages
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading message:", err)
			continue
		}

		// Parse event
		var event GameEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("Error parsing event: %v", err)
			continue
		}

		// Log the event
		log.Printf("Game Event: ID=%s Winner=%s Duration=%ds",
			event.GameID, event.Winner, event.Duration)

		// Process analytics
		if err := processGameEvent(event); err != nil {
			log.Printf("Error processing analytics: %v", err)
		}
	}
}

func processGameEvent(event GameEvent) error {
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
		return err
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

	return nil
}
