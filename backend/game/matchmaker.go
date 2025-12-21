package game

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Matchmaker struct {
	Queue []*Player
	Mutex sync.Mutex
}

var GlobalMatchmaker = &Matchmaker{
	Queue: make([]*Player, 0),
}

func (m *Matchmaker) IsPlayerInQueue(username string) bool {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	for _, p := range m.Queue {
		if p.Username == username {
			return true
		}
	}
	return false
}

func (m *Matchmaker) RemovePlayer(player *Player) bool {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	for i, p := range m.Queue {
		if p == player || p.ID == player.ID {
			// Remove player from queue
			m.Queue = append(m.Queue[:i], m.Queue[i+1:]...)
			log.Printf("Player %s removed from queue. Queue size: %d", player.Username, len(m.Queue))
			return true
		}
	}
	return false
}

func (m *Matchmaker) AddPlayer(p *Player) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	m.Queue = append(m.Queue, p)
	log.Printf("Player %s added to queue. Queue size: %d", p.Username, len(m.Queue))

	if len(m.Queue) >= 2 {
		p1 := m.Queue[0]
		p2 := m.Queue[1]
		m.Queue = m.Queue[2:]
		m.StartGame(p1, p2)
	} else {
		go m.WaitForMatch(p)
	}
}

func (m *Matchmaker) WaitForMatch(p *Player) {
	time.Sleep(10 * time.Second)

	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	found := false
	for i, qp := range m.Queue {
		if qp == p {
			m.Queue = append(m.Queue[:i], m.Queue[i+1:]...)
			found = true
			break
		}
	}

	if found {
		log.Printf("Timeout for player %s. Starting bot game.", p.Username)
		bot := &Player{
			ID:       "bot-" + uuid.New().String(),
			Username: "Bot",
			IsBot:    true,
		}
		m.StartGame(p, bot)
	}
}

func (m *Matchmaker) StartGame(p1, p2 *Player) {
	gameID := uuid.New().String()
	game := NewGame(gameID, p1, p2)

	go game.Start()

	GameManagerInstance.AddGame(game)
}

// GameManager to keep track of active games
type GameManager struct {
	Games       map[string]*Game
	PlayerGames map[string]string // PlayerID -> GameID
	Mutex       sync.RWMutex
}

var GameManagerInstance = &GameManager{
	Games:       make(map[string]*Game),
	PlayerGames: make(map[string]string),
}

func (gm *GameManager) AddGame(g *Game) {
	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()
	gm.Games[g.ID] = g
	gm.PlayerGames[g.Player1.ID] = g.ID
	if !g.Player2.IsBot {
		gm.PlayerGames[g.Player2.ID] = g.ID
	}
}

func (gm *GameManager) GetGame(id string) *Game {
	gm.Mutex.RLock()
	defer gm.Mutex.RUnlock()
	return gm.Games[id]
}

func (gm *GameManager) GetGameByPlayerID(playerID string) *Game {
	gm.Mutex.RLock()
	defer gm.Mutex.RUnlock()
	gameID, ok := gm.PlayerGames[playerID]
	if !ok {
		return nil
	}
	return gm.Games[gameID]
}

func (gm *GameManager) RemoveGame(id string) {
	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()

	if g, ok := gm.Games[id]; ok {
		delete(gm.PlayerGames, g.Player1.ID)
		if !g.Player2.IsBot {
			delete(gm.PlayerGames, g.Player2.ID)
		}
		delete(gm.Games, id)
	}
}

func (gm *GameManager) IsUsernameTaken(username string) bool {
	gm.Mutex.RLock()
	defer gm.Mutex.RUnlock()

	for _, g := range gm.Games {
		if g.State == "active" {
			if g.Player1.Username == username && g.Player1.IsConnected {
				return true
			}
			if !g.Player2.IsBot && g.Player2.Username == username && g.Player2.IsConnected {
				return true
			}
		}
	}
	return false
}

func (gm *GameManager) GetPlayerByUsername(username string) (*Player, *Game) {
	gm.Mutex.RLock()
	defer gm.Mutex.RUnlock()

	for _, g := range gm.Games {
		if g.State == "active" {
			if g.Player1.Username == username {
				return g.Player1, g
			}
			if !g.Player2.IsBot && g.Player2.Username == username {
				return g.Player2, g
			}
		}
	}
	return nil, nil
}
