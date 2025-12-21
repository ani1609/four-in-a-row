package main

import (
	"log"
	"net/http"

	"4-in-a-row/analytics"
	"4-in-a-row/config"
	"4-in-a-row/db"
	"4-in-a-row/handlers"
)

func main() {
	log.Println("Starting 4-in-a-row Game Server...")

	cfg := config.Load()

	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	streamConfig := make(map[string]string)

	switch cfg.EventStream {
	case "kafka":
		if cfg.KafkaBrokers != "" {
			streamConfig["brokers"] = cfg.KafkaBrokers
			streamConfig["topic"] = cfg.KafkaTopic
			if err := analytics.InitEventStream("kafka", streamConfig); err != nil {
				log.Printf("Failed to initialize Kafka: %v", err)
			} else {
				log.Println("Kafka analytics enabled")
			}
		} else {
			log.Println("Kafka analytics disabled (no KAFKA_BROKERS_LOCAL set)")
		}
	case "redis":
		if cfg.RedisURL != "" {
			streamConfig["url"] = cfg.RedisURL
			streamConfig["password"] = cfg.RedisPassword
			streamConfig["stream"] = cfg.RedisStream
			if err := analytics.InitEventStream("redis", streamConfig); err != nil {
				log.Printf("Failed to initialize Redis: %v", err)
			} else {
				log.Println("Redis Streams analytics enabled")
			}
		} else {
			log.Println("Redis analytics disabled (no REDIS_URL set)")
		}
	}

	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/ws", handlers.WSHandler)
	http.HandleFunc("/leaderboard", handlers.LeaderboardHandler)
	http.HandleFunc("/metrics", handlers.GameMetricsHandler)
	http.HandleFunc("/recent-games", handlers.RecentGamesHandler)

	log.Printf("Server listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatal("Server error: ", err)
	}
}
