package handlers

import (
	"log"
	"net/http"

	"4-in-a-row/game"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade WS:", err)
		return
	}
	defer conn.Close()

	log.Println("New Client Connected")

	var currentPlayer *game.Player

	defer func() {
		if currentPlayer != nil {
			removed := game.GlobalMatchmaker.RemovePlayer(currentPlayer)
			if removed {
				log.Printf("Player %s removed from matchmaking queue on disconnect", currentPlayer.Username)
			}

			g := game.GameManagerInstance.GetGameByPlayerID(currentPlayer.ID)
			if g != nil {
				g.HandleDisconnect(currentPlayer)
			}
		}
	}()

	for {
		var msg game.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WS Read Error: %v", err)
			}
			break
		}

		switch msg.Type {
		case game.MsgJoinQueue:
			username, ok := msg.Payload.(string)
			if !ok {
				username = "Anonymous"
			}

			if game.GlobalMatchmaker.IsPlayerInQueue(username) {
				conn.WriteJSON(game.Message{Type: game.MsgError, Payload: "Username already taken"})
				continue
			}

			existingPlayer, existingGame := game.GameManagerInstance.GetPlayerByUsername(username)
			if existingPlayer != nil && existingGame != nil {
				if !existingPlayer.IsConnected {
					log.Printf("User %s reconnecting to game %s", username, existingGame.ID)
					reconnectedPlayer, success := existingGame.HandleReconnect(existingPlayer.ID, conn)
					if success {
						currentPlayer = reconnectedPlayer
					} else {
						conn.WriteJSON(game.Message{Type: game.MsgError, Payload: "Reconnect failed"})
					}
					continue
				} else {
					conn.WriteJSON(game.Message{Type: game.MsgError, Payload: "Username already taken"})
					continue
				}
			}

			player := &game.Player{
				ID:       uuid.New().String(),
				Username: username,
				Conn:     conn,
			}

			currentPlayer = player
			game.GlobalMatchmaker.AddPlayer(player)

		case game.MsgReconnect:
			playerID, ok := msg.Payload.(string)
			if !ok {
				continue
			}

			g := game.GameManagerInstance.GetGameByPlayerID(playerID)
			if g != nil {
				p, success := g.HandleReconnect(playerID, conn)
				if success {
					currentPlayer = p
				} else {
					conn.WriteJSON(game.Message{Type: game.MsgError, Payload: "Reconnect failed or game ended"})
				}
			} else {
				conn.WriteJSON(game.Message{Type: game.MsgError, Payload: "Game not found"})
			}

		case game.MsgMove:
			payload, ok := msg.Payload.(map[string]interface{})
			if !ok {
				continue
			}

			gameID, _ := payload["gameId"].(string)
			col, _ := payload["column"].(float64)

			g := game.GameManagerInstance.GetGame(gameID)
			if g != nil {
				var p *game.Player
				if currentPlayer != nil {
					p = currentPlayer
				} else if g.Player1.Conn == conn {
					p = g.Player1
				} else if g.Player2.Conn == conn {
					p = g.Player2
				}

				if p != nil {
					g.HandleMove(p, int(col))
				}
			}
		}
	}
}
