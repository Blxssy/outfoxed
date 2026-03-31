package ws

import (
	"sync"
)

type Room struct {
	id string

	mu    sync.RWMutex
	conns map[*Conn]struct{}
}

func NewRoom(gameID string) *Room {
	return &Room{
		id:    gameID,
		conns: make(map[*Conn]struct{}),
	}
}

func (r *Room) Add(c *Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.conns[c] = struct{}{}
}

func (r *Room) Remove(c *Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.conns, c)
}

func (r *Room) Broadcast(v any) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for c := range r.conns {
		c.SendJSON(v)
	}
}

func (r *Room) Size() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.conns)
}
