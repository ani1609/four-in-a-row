package game

import (
	"time"
	"github.com/gorilla/websocket"
)

type Player struct {
	ID       string
	Username string
	Conn     *websocket.Conn
	IsBot    bool
	Symbol         int // 1 or 2
	IsConnected    bool
	DisconnectedAt time.Time
}

func (p *Player) SendMessage(msg interface{}) error {
	if p.IsBot || p.Conn == nil {
		return nil
	}
	return p.Conn.WriteJSON(msg)
}
