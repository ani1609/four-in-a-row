package db

import (
	"4-in-a-row/config"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PlayerData struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Symbol   int    `json:"symbol"`
	Type     string `json:"type"` // "human" or "bot"
}

type MoveData struct {
	MoveNumber int   `json:"moveNumber"`
	Player     int   `json:"player"`
	Column     int   `json:"column"`
	Row        int   `json:"row"`
	Timestamp  int64 `json:"timestamp"` // Unix timestamp
}

type GameResult struct {
	ID        uint       `gorm:"primaryKey"`
	GameID    string     `gorm:"index"`
	Player1   PlayerData `gorm:"type:jsonb;serializer:json"`
	Player2   PlayerData `gorm:"type:jsonb;serializer:json"`
	Winner    string     `gorm:"index"`
	Moves     []MoveData `gorm:"type:jsonb;serializer:json"`
	Duration  int64
	CreatedAt time.Time
}

var DB *gorm.DB

// InitDB initializes the database connection and runs migrations
func InitDB() error {
	cfg := config.Get()
	dsn := cfg.DatabaseURL
	if dsn == "" {
		return &ConfigError{"DATABASE_URL not configured"}
	}

	log.Println("Connecting to database...")

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto-migrate schema
	if err := DB.AutoMigrate(&GameResult{}, &PlayerStats{}, &GameMetrics{}); err != nil {
		return err
	}

	log.Println("Database connected successfully")
	return nil
}

// SaveGameResult persists a completed game to the database
func SaveGameResult(gameID string, p1, p2 PlayerData, winner string, moves []MoveData, duration int64) {
	if DB == nil {
		log.Println("Database not initialized, skipping save")
		return
	}

	result := GameResult{
		GameID:    gameID,
		Player1:   p1,
		Player2:   p2,
		Winner:    winner,
		Moves:     moves,
		Duration:  duration,
		CreatedAt: time.Now(),
	}

	if err := DB.Create(&result).Error; err != nil {
		log.Printf("Failed to save game result: %v", err)
	} else {
		log.Printf("Game result saved: %s won, %d moves recorded", winner, len(moves))
	}
}

// ConfigError represents a configuration error
type ConfigError struct {
	msg string
}

func (e *ConfigError) Error() string {
	return e.msg
}
