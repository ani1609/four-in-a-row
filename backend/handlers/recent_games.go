package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"4-in-a-row/db"
)

type RecentGameResponse struct {
	GameID     string `json:"gameId"`
	Player1    string `json:"player1"`
	Player2    string `json:"player2"`
	Winner     string `json:"winner"`
	Duration   int64  `json:"duration"`
	TotalMoves int    `json:"totalMoves"`
	PlayedAt   string `json:"playedAt"`
}

func RecentGamesHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if db.DB == nil {
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	// Get recent games (last 20)
	var games []db.GameResult
	err := db.DB.Order("created_at DESC").Limit(20).Find(&games).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Transform to response format
	response := make([]RecentGameResponse, len(games))
	for i, game := range games {
		response[i] = RecentGameResponse{
			GameID:     game.GameID,
			Player1:    game.Player1.Username,
			Player2:    game.Player2.Username,
			Winner:     game.Winner,
			Duration:   game.Duration,
			TotalMoves: len(game.Moves),
			PlayedAt:   game.CreatedAt.Format(time.RFC3339),
		}
	}

	json.NewEncoder(w).Encode(response)
}
