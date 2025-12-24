package game

// Message Types
const (
	MsgJoinQueue    = "JOIN_QUEUE"
	MsgGameStart    = "GAME_START"
	MsgMove         = "MOVE"
	MsgUpdate       = "GAME_UPDATE"
	MsgGameOver     = "GAME_OVER"
	MsgError        = "ERROR"
	MsgReconnect    = "RECONNECT"
	MsgPlayerStatus = "PLAYER_STATUS"
)

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type PlayerInfo struct {
	PlayerID string `json:"playerId,omitempty"` // Only for "you"
	Username string `json:"username"`
	Symbol   int    `json:"symbol"`
	Type     string `json:"type"` // "human" or "bot"
	IsOnline bool   `json:"isOnline"`
}

type GameStartPayload struct {
	GameID   string     `json:"gameId"`
	You      PlayerInfo `json:"you"`
	Opponent PlayerInfo `json:"opponent"`
	YourTurn bool       `json:"yourTurn"`
}

type LastMove struct {
	Player int `json:"player"`
	Column int `json:"column"`
	Row    int `json:"row"`
}

type GameUpdatePayload struct {
	Grid        [Rows][Cols]int `json:"grid"`
	CurrentTurn int             `json:"currentTurn"` // 1 or 2
	LastMove    *LastMove       `json:"lastMove,omitempty"`
	MoveNumber  int             `json:"moveNumber"`
}

type ReconnectPayload struct {
	GameID      string          `json:"gameId"`
	You         PlayerInfo      `json:"you"`
	Opponent    PlayerInfo      `json:"opponent"`
	Grid        [Rows][Cols]int `json:"grid"`
	CurrentTurn int             `json:"currentTurn"`
	YourTurn    bool            `json:"yourTurn"`
	MoveNumber  int             `json:"moveNumber"`
}

type GameOverPayload struct {
	Winner string `json:"winner"` // "1", "2", or "draw"
}

type PlayerStatusPayload struct {
	PlayerSymbol int  `json:"playerSymbol"`
	IsOnline     bool `json:"isOnline"`
	TimeLeft     int  `json:"timeLeft"` // Seconds left before forfeit (0 if online)
}
