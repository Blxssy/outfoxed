package ws

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"fox/internal/modules/game/domain"
	"fox/internal/modules/game/service"

	cws "github.com/coder/websocket"
	"github.com/go-chi/chi/v5"
)

type Auth interface {
	ParseToken(token string) (string, error)
}

type GameStateGetter interface {
	GetState(ctx context.Context, gameID string, userID string) (domain.GameState, error)
}

type Handler struct {
	hub   *Hub
	auth  Auth
	game  *service.Service
	state GameStateGetter
}

func NewHandler(hub *Hub, auth Auth, game *service.Service, state GameStateGetter) *Handler {
	return &Handler{hub: hub, auth: auth, game: game, state: state}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "id")
	if gameID == "" {
		http.Error(w, "missing game id", http.StatusBadRequest)
		return
	}

	userID, err := h.userIDFromRequest(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	wsConn, err := cws.Accept(w, r, &cws.AcceptOptions{
		OriginPatterns: []string{"http://localhost:*", "https://your-domain.com"},
	})
	if err != nil {
		return
	}
	defer func() { _ = wsConn.Close(cws.StatusNormalClosure, "bye") }()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	conn := NewConn(wsConn)

	room := h.hub.GetRoom(gameID)
	room.Add(conn)
	defer func() {
		room.Remove(conn)
		conn.Close()
		h.hub.RemoveRoomIfEmpty(gameID)
	}()

	go conn.WriteLoop(ctx)

	st, err := h.state.GetState(ctx, gameID, userID)
	if err != nil {
		conn.SendJSON(service.NewErrorResponse("", "forbidden", "You do not have access to this game."))
		return
	}
	conn.SendJSON(service.NewUpdateResponse("", st, nil))

	for {
		readCtx, cancelRead := context.WithTimeout(ctx, 60*time.Second)
		msgType, data, err := wsConn.Read(readCtx)
		cancelRead()

		if err != nil {
			return
		}
		if msgType != cws.MessageText {
			continue
		}

		var req service.WSRequest
		if err := json.Unmarshal(data, &req); err != nil {
			conn.SendJSON(service.NewErrorResponse("", "bad_json", "Invalid JSON."))
			continue
		}

		actorID := domain.PlayerID(userID)
		cmd, err := service.DecodeCommand(req, actorID)
		if err != nil {
			conn.SendJSON(service.NewErrorResponse(req.ID, "invalid_command", err.Error()))
			continue
		}

		newState, events, err := h.game.ApplyCommand(ctx, gameID, userID, cmd)
		if err != nil {
			conn.SendJSON(service.ErrorResponse(req.ID, err))
			continue
		}

		room.Broadcast(service.NewUpdateResponse(req.ID, newState, events))
	}
}

func (h *Handler) userIDFromRequest(r *http.Request) (string, error) {
	// 1) Bearer header
	if ah := r.Header.Get("Authorization"); ah != "" {
		parts := strings.SplitN(ah, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") && parts[1] != "" {
			return h.auth.ParseToken(parts[1])
		}
		return "", errors.New("bad authorization header")
	}

	// 2) Query token
	if token := r.URL.Query().Get("token"); token != "" {
		return h.auth.ParseToken(token)
	}

	return "", errors.New("missing authorization")
}
