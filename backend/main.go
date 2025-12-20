package main

import (
	"log"
	"net/http"
	"os"

	"4-in-a-row/analytics"
	"4-in-a-row/db"
	"4-in-a-row/handlers"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting 4-in-a-row Game Server...")

	// Load .env file for local development (ignored in production)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize Database
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Kafka (optional)
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers != "" {
		analytics.InitKafka([]string{kafkaBrokers}, "game-analytics")
		log.Println("Kafka analytics enabled")
	} else {
		log.Println("Kafka analytics disabled")
	}

	// Setup HTTP routes
	http.HandleFunc("/ws", handlers.WSHandler)
	http.HandleFunc("/leaderboard", handlers.LeaderboardHandler)
	http.HandleFunc("/metrics", handlers.GameMetricsHandler)
	http.HandleFunc("/recent-games", handlers.RecentGamesHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server error: ", err)
	}
}
