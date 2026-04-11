package ws

import (
	"sync"
)

type Hub struct {
	mu    sync.RWMutex
	rooms map[string]*Room
}

func NewHub() *Hub {
	return &Hub{rooms: make(map[string]*Room)}
}

func (h *Hub) GetRoom(gameID string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	r, ok := h.rooms[gameID]
	if !ok {
		r = NewRoom(gameID)
		h.rooms[gameID] = r
	}
	return r
}

func (h *Hub) RemoveRoomIfEmpty(gameID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	r, ok := h.rooms[gameID]
	if !ok {
		return
	}
	if r.Size() == 0 {
		delete(h.rooms, gameID)
	}
}
