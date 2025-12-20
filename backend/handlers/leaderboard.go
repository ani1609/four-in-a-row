package handlers

import (
	"encoding/json"
	"net/http"

	"4-in-a-row/db"
)

type LeaderboardEntry struct {
	Winner string `json:"name"`
	Wins   int64  `json:"wins"`
}

func LeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	if db.DB == nil {
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	var results []LeaderboardEntry
	err := db.DB.Model(&db.GameResult{}).
		Select("winner, count(*) as wins").
		Where("winner != ?", "draw").
		Group("winner").
		Order("wins desc").
		Limit(10).
		Scan(&results).Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(results)
}
