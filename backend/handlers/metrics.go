package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"4-in-a-row/db"
)

// GameMetricsResponse represents the overall game analytics
type GameMetricsResponse struct {
	TotalGames      int              `json:"totalGames"`
	TotalPlayers    int              `json:"totalPlayers"`
	AverageDuration float64          `json:"averageDuration"`
	GamesToday      int              `json:"gamesToday"`
	RecentActivity  []HourlyActivity `json:"recentActivity"`
}

type HourlyActivity struct {
	Hour            int     `json:"hour"`
	GamesPlayed     int     `json:"gamesPlayed"`
	AverageDuration float64 `json:"averageDuration"`
}

func GameMetricsHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if db.DB == nil {
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	// Get total games
	var totalGames int64
	db.DB.Model(&db.GameResult{}).Count(&totalGames)

	// Get total unique players
	var totalPlayers int64
	db.DB.Model(&db.PlayerStats{}).Count(&totalPlayers)

	// Get average duration (use COALESCE to handle NULL when no games exist)
	var avgDuration float64
	db.DB.Model(&db.GameResult{}).Select("COALESCE(AVG(duration), 0)").Scan(&avgDuration)

	// Get games today
	today := time.Now().Truncate(24 * time.Hour)
	var gamesToday int64
	db.DB.Model(&db.GameResult{}).Where("created_at >= ?", today).Count(&gamesToday)

	// Get recent hourly activity (last 24 hours)
	var recentMetrics []db.GameMetrics
	yesterday := time.Now().Add(-24 * time.Hour)
	db.DB.Where("date >= ?", yesterday.Truncate(24*time.Hour)).
		Order("date DESC, hour DESC").
		Limit(24).
		Find(&recentMetrics)

	// Convert to response format
	recentActivity := make([]HourlyActivity, 0, len(recentMetrics))
	for _, m := range recentMetrics {
		recentActivity = append(recentActivity, HourlyActivity{
			Hour:            m.Hour,
			GamesPlayed:     m.GamesPlayed,
			AverageDuration: m.AverageDuration,
		})
	}

	response := GameMetricsResponse{
		TotalGames:      int(totalGames),
		TotalPlayers:    int(totalPlayers),
		AverageDuration: avgDuration,
		GamesToday:      int(gamesToday),
		RecentActivity:  recentActivity,
	}

	json.NewEncoder(w).Encode(response)
}
