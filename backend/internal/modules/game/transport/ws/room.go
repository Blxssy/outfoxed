package ws

import "sync"

type Client struct {
	UserID string
	Conn   *Conn
}

type Room struct {
	id string

	mu      sync.RWMutex
	clients map[*Conn]Client
}

func NewRoom(gameID string) *Room {
	return &Room{
		id:      gameID,
		clients: make(map[*Conn]Client),
	}
}

func (r *Room) Add(userID string, c *Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.clients[c] = Client{
		UserID: userID,
		Conn:   c,
	}
}

func (r *Room) Remove(c *Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, c)
}

// Clients возвращает снимок клиентов комнаты.
func (r *Room) Clients() []Client {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]Client, 0, len(r.clients))
	for _, client := range r.clients {
		out = append(out, client)
	}
	return out
}

func (r *Room) Size() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.clients)
}

func (r *Room) ConnectedUserIDs() map[string]bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make(map[string]bool, len(r.clients))
	for _, client := range r.clients {
		out[client.UserID] = true
	}
	return out
}
