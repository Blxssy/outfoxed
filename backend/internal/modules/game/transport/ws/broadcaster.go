package ws

import (
	"context"

	"fox/internal/modules/game/domain"
	"fox/internal/modules/game/service"
)

type LobbySnapshotProvider interface {
	GetLobby(ctx context.Context, gameID string, userID string) (service.LobbySnapshot, error)
}

type Broadcaster struct {
	hub           *Hub
	lobbyProvider LobbySnapshotProvider
}

func NewBroadcaster(hub *Hub) *Broadcaster {
	return &Broadcaster{
		hub: hub,
	}
}

func (b *Broadcaster) SetLobbyProvider(provider LobbySnapshotProvider) {
	b.lobbyProvider = provider
}

func (b *Broadcaster) PublishLobby(gameID string) {
	if b.lobbyProvider == nil {
		return
	}

	room, ok := b.hub.FindRoom(gameID)
	if !ok || room.Size() == 0 {
		return
	}

	clients := room.Clients()
	ctx := context.Background()

	for _, client := range clients {
		snapshot, err := b.lobbyProvider.GetLobby(ctx, gameID, client.UserID)
		if err != nil {
			continue
		}

		client.Conn.SendJSON(service.NewLobbyUpdateResponse("", snapshot))
	}
}

func (b *Broadcaster) PublishGame(gameID string, state domain.GameState, events []domain.Event) {
	room, ok := b.hub.FindRoom(gameID)
	if !ok || room.Size() == 0 {
		return
	}

	clients := room.Clients()
	connected := room.ConnectedUserIDs()

	for _, client := range clients {
		view := domain.BuildGameView(state, domain.PlayerID(client.UserID))
		applyConnectedState(&view, connected)
		client.Conn.SendJSON(service.NewGameUpdateResponse("", view, events))
	}
}
