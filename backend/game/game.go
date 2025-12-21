package game

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"4-in-a-row/analytics"
	"4-in-a-row/db"
)

type Game struct {
	ID         string
	Player1    *Player
	Player2    *Player
	Board      *Board
	Turn       int    // 1 or 2
	State      string // "active", "finished"
	Winner     int    // 0 = none, 1 = p1, 2 = p2, 3 = draw
	Mutex      sync.Mutex
	LastMove   time.Time
	Moves      []db.MoveData
	MoveNumber int
	StartTime  time.Time
}

func NewGame(id string, p1, p2 *Player) *Game {
	p1.Symbol = 1
	p2.Symbol = 2
	return &Game{
		ID:         id,
		Player1:    p1,
		Player2:    p2,
		Board:      NewBoard(),
		Turn:       1, // Player 1 starts
		State:      "active",
		LastMove:   time.Now(),
		Moves:      []db.MoveData{},
		MoveNumber: 0,
		StartTime:  time.Now(),
	}
}

func (g *Game) Start() {
	g.Player1.SendMessage(Message{
		Type: MsgGameStart,
		Payload: GameStartPayload{
			GameID: g.ID,
			You: PlayerInfo{
				PlayerID: g.Player1.ID,
				Username: g.Player1.Username,
				Symbol:   1,
				Type:     getPlayerType(g.Player1),
			},
			Opponent: PlayerInfo{
				Username: g.Player2.Username,
				Symbol:   2,
				Type:     getPlayerType(g.Player2),
			},
			YourTurn: true,
		},
	})

	g.Player1.IsConnected = true

	g.Player2.SendMessage(Message{
		Type: MsgGameStart,
		Payload: GameStartPayload{
			GameID: g.ID,
			You: PlayerInfo{
				PlayerID: g.Player2.ID,
				Username: g.Player2.Username,
				Symbol:   2,
				Type:     getPlayerType(g.Player2),
			},
			Opponent: PlayerInfo{
				Username: g.Player1.Username,
				Symbol:   1,
				Type:     getPlayerType(g.Player1),
			},
			YourTurn: false,
		},
	})

	g.Player2.IsConnected = true

	log.Printf("Game %s started: %s vs %s", g.ID, g.Player1.Username, g.Player2.Username)
}

func (g *Game) HandleMove(player *Player, col int) {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	if g.State != "active" {
		return
	}

	if player.Symbol != g.Turn {
		player.SendMessage(Message{Type: MsgError, Payload: "Not your turn"})
		return
	}

	row, err := g.Board.DropDisc(col, g.Turn)
	if err != nil {
		player.SendMessage(Message{Type: MsgError, Payload: "Invalid move"})
		return
	}

	g.MoveNumber++
	moveData := db.MoveData{
		MoveNumber: g.MoveNumber,
		Player:     g.Turn,
		Column:     col,
		Row:        row,
		Timestamp:  time.Now().Unix(),
	}
	g.Moves = append(g.Moves, moveData)

	if g.Board.CheckWin(row, col, g.Turn) {
		g.State = "finished"
		g.Winner = g.Turn
		g.BroadcastUpdate(row, col)
		g.BroadcastGameOver()
		return
	}

	if g.Board.IsFull() {
		g.State = "finished"
		g.Winner = 3 // Draw
		g.BroadcastUpdate(row, col)
		g.BroadcastGameOver()
		return
	}

	if g.Turn == 1 {
		g.Turn = 2
	} else {
		g.Turn = 1
	}

	g.BroadcastUpdate(row, col)

	if g.Turn == 2 && g.Player2.IsBot {
		go g.TriggerBotMove()
	}
}

func (g *Game) BroadcastUpdate(lastRow, lastCol int) {
	var lastMove *LastMove
	if g.MoveNumber > 0 {
		lastMoveData := g.Moves[len(g.Moves)-1]
		lastMove = &LastMove{
			Player: lastMoveData.Player,
			Column: lastMoveData.Column,
			Row:    lastMoveData.Row,
		}
	}

	msg := Message{
		Type: MsgUpdate,
		Payload: GameUpdatePayload{
			Grid:        g.Board.Grid,
			CurrentTurn: g.Turn,
			LastMove:    lastMove,
			MoveNumber:  g.MoveNumber,
		},
	}
	g.Player1.SendMessage(msg)
	g.Player2.SendMessage(msg)
}

func (g *Game) BroadcastGameOver() {
	winnerStr := ""
	winnerName := ""

	if g.Winner == 1 {
		winnerStr = g.Player1.Username
		winnerName = g.Player1.Username
	} else if g.Winner == 2 {
		winnerStr = g.Player2.Username
		winnerName = g.Player2.Username
	} else {
		winnerStr = "draw"
		winnerName = "Draw"
	}

	msg := Message{
		Type: MsgGameOver,
		Payload: GameOverPayload{
			Winner: winnerStr,
		},
	}
	g.Player1.SendMessage(msg)
	g.Player2.SendMessage(msg)

	duration := int64(time.Since(g.StartTime).Seconds())

	// Log game result
	if g.Winner == 3 {
		log.Printf("Game %s ended in a DRAW between %s and %s (Duration: %ds)",
			g.ID, g.Player1.Username, g.Player2.Username, duration)
	} else {
		log.Printf("Game %s WON by %s! (%s vs %s, Duration: %ds)",
			g.ID, winnerName, g.Player1.Username, g.Player2.Username, duration)
	}

	// Prepare player data
	p1Data := db.PlayerData{
		ID:       g.Player1.ID,
		Username: g.Player1.Username,
		Symbol:   g.Player1.Symbol,
		Type:     getPlayerType(g.Player1),
	}
	p2Data := db.PlayerData{
		ID:       g.Player2.ID,
		Username: g.Player2.Username,
		Symbol:   g.Player2.Symbol,
		Type:     getPlayerType(g.Player2),
	}

	// Persist game result with moves
	db.SaveGameResult(g.ID, p1Data, p2Data, winnerStr, g.Moves, duration)

	// Emit Kafka event
	analytics.EmitGameEnd(g.ID, winnerStr, duration)
}

func (g *Game) TriggerBotMove() {
	// Realistic delay
	time.Sleep(500 * time.Millisecond)

	// Use smart bot AI
	bot := &BotAI{
		board:          g.Board,
		botSymbol:      2, // Bot is always player 2
		opponentSymbol: 1,
	}

	col := bot.GetBestMove()

	if col != -1 {
		g.HandleMove(g.Player2, col)
	}
}

func (g *Game) HandleDisconnect(player *Player) {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	if g.State != "active" {
		return
	}

	log.Printf("Player %s disconnected from game %s", player.Username, g.ID)
	player.IsConnected = false
	player.DisconnectedAt = time.Now()

	go func() {
		time.Sleep(30 * time.Second)
		g.Mutex.Lock()
		defer g.Mutex.Unlock()

		if g.State == "active" && !player.IsConnected {
			log.Printf("Player %s timed out. Forfeiting game %s.", player.Username, g.ID)

			g.State = "finished"

			if player.Symbol == 1 {
				g.Winner = 2
			} else {
				g.Winner = 1
			}

			g.BroadcastGameOver()

			GameManagerInstance.RemoveGame(g.ID)
		}
	}()
}

func (g *Game) HandleReconnect(playerID string, conn *websocket.Conn) (*Player, bool) {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	if g.State != "active" && g.State != "finished" {
		return nil, false
	}

	var p *Player
	if g.Player1.ID == playerID {
		p = g.Player1
	} else if g.Player2.ID == playerID {
		p = g.Player2
	} else {
		return nil, false
	}

	p.Conn = conn
	p.IsConnected = true

	log.Printf("Player %s reconnected to game %s", p.Username, g.ID)

	opponent := GetOpponent(g, p)

	reconnectMsg := Message{
		Type: MsgReconnect,
		Payload: ReconnectPayload{
			GameID: g.ID,
			You: PlayerInfo{
				PlayerID: p.ID,
				Username: p.Username,
				Symbol:   p.Symbol,
				Type:     getPlayerType(p),
			},
			Opponent: PlayerInfo{
				Username: opponent.Username,
				Symbol:   opponent.Symbol,
				Type:     getPlayerType(opponent),
			},
			Grid:        g.Board.Grid,
			CurrentTurn: g.Turn,
			YourTurn:    g.Turn == p.Symbol && g.State == "active",
			MoveNumber:  g.MoveNumber,
		},
	}
	p.SendMessage(reconnectMsg)

	if g.State == "finished" {
		winnerStr := ""
		if g.Winner == 1 {
			winnerStr = g.Player1.Username
		} else if g.Winner == 2 {
			winnerStr = g.Player2.Username
		} else {
			winnerStr = "draw"
		}

		p.SendMessage(Message{
			Type: MsgGameOver,
			Payload: GameOverPayload{
				Winner: winnerStr,
			},
		})
	}

	return p, true
}

func GetOpponent(g *Game, p *Player) *Player {
	if p == g.Player1 {
		return g.Player2
	}
	return g.Player1
}

func getPlayerType(p *Player) string {
	if p.IsBot {
		return "bot"
	}
	return "human"
}
