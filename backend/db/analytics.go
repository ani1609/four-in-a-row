package db

import (
	"time"
)

// PlayerStats tracks user-specific metrics
type PlayerStats struct {
	ID            uint   `gorm:"primaryKey"`
	Username      string `gorm:"uniqueIndex"`
	TotalGames    int
	Wins          int
	Losses        int
	Draws         int
	TotalDuration int64 // Total seconds played
	LastPlayed    time.Time
	UpdatedAt     time.Time
}

// GameMetrics tracks overall gameplay metrics
type GameMetrics struct {
	ID              uint      `gorm:"primaryKey"`
	Date            time.Time `gorm:"index"`
	Hour            int
	GamesPlayed     int
	TotalDuration   int64
	AverageDuration float64
	UpdatedAt       time.Time
}

// UpdatePlayerStats updates player statistics for a completed game
func UpdatePlayerStats(username string, won bool, isDraw bool, duration int64) error {
	if DB == nil {
		return &ConfigError{"Database not initialized"}
	}

	var stats PlayerStats

	result := DB.Where(PlayerStats{Username: username}).FirstOrCreate(&stats)
	if result.Error != nil {
		return result.Error
	}

	stats.TotalGames++
	stats.TotalDuration += duration
	stats.LastPlayed = time.Now()

	if isDraw {
		stats.Draws++
	} else if won {
		stats.Wins++
	} else {
		stats.Losses++
	}

	return DB.Save(&stats).Error
}

// UpdateGameMetrics aggregates game metrics by hour
func UpdateGameMetrics(duration int64, timestamp time.Time) error {
	if DB == nil {
		return &ConfigError{"Database not initialized"}
	}

	date := time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, timestamp.Location())
	hour := timestamp.Hour()

	var metrics GameMetrics

	result := DB.Where(GameMetrics{Date: date, Hour: hour}).FirstOrCreate(&metrics)
	if result.Error != nil {
		return result.Error
	}

	metrics.GamesPlayed++
	metrics.TotalDuration += duration
	metrics.AverageDuration = float64(metrics.TotalDuration) / float64(metrics.GamesPlayed)

	return DB.Save(&metrics).Error
}

// GetTopPlayers returns top N players by wins
func GetTopPlayers(limit int) ([]PlayerStats, error) {
	if DB == nil {
		return nil, &ConfigError{"Database not initialized"}
	}

	var players []PlayerStats
	err := DB.Order("wins DESC").Limit(limit).Find(&players).Error
	return players, err
}
