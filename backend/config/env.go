package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ResourceEnvironment string
	EventStream         string
	Port                string
	DatabaseURL         string
	RedisURL            string
	RedisPassword       string
	RedisStream         string
	KafkaBrokers        string
	KafkaTopic          string
}

var globalConfig *Config

func Load() *Config {
	if globalConfig != nil {
		return globalConfig
	}

	godotenv.Load()

	resourceEnv := strings.ToLower(os.Getenv("RESOURCE_ENVIRONMENT"))
	if resourceEnv == "" {
		resourceEnv = "local"
	}

	config := &Config{
		ResourceEnvironment: resourceEnv,
		EventStream:         getEnv("EVENT_STREAM", "kafka"),
		Port:                getEnv("PORT", "8080"),
		RedisStream:         getEnv("REDIS_STREAM", "game-events"),
	}

	if resourceEnv == "cloud" {
		config.DatabaseURL = os.Getenv("DATABASE_URL_CLOUD")
		if config.DatabaseURL == "" {
			log.Fatal("DATABASE_URL_CLOUD not set for cloud environment")
		}
	} else {
		config.DatabaseURL = os.Getenv("DATABASE_URL_LOCAL")
		if config.DatabaseURL == "" {
			log.Fatal("DATABASE_URL_LOCAL not set for local environment")
		}
	}

	if resourceEnv == "cloud" {
		config.RedisURL = os.Getenv("REDIS_URL_CLOUD")
		config.RedisPassword = ""
	} else {
		config.RedisURL = getEnv("REDIS_URL_LOCAL", "localhost:6379")
		config.RedisPassword = os.Getenv("REDIS_PASSWORD_LOCAL")
	}

	config.KafkaBrokers = getEnv("KAFKA_BROKERS_LOCAL", "localhost:9092")
	config.KafkaTopic = getEnv("KAFKA_TOPIC_LOCAL", "game-analytics")

	log.Printf("Configuration loaded: environment=%s, event_stream=%s",
		config.ResourceEnvironment, config.EventStream)

	globalConfig = config
	return config
}

func Get() *Config {
	if globalConfig == nil {
		return Load()
	}
	return globalConfig
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
