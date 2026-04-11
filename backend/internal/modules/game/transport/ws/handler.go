package ws

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	authService "fox/internal/modules/auth/service"
	"fox/internal/modules/game/domain"
	"fox/internal/modules/game/service"

	cws "github.com/coder/websocket"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type Handler struct {
	log          zerolog.Logger
	hub          *Hub
	gameService  *service.Service
	tokenManager *authService.TokenManager
}

func NewHandler(log zerolog.Logger, hub *Hub, gameService *service.Service, tokenManager *authService.TokenManager) *Handler {
	return &Handler{log: log, hub: hub, gameService: gameService, tokenManager: tokenManager}
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

	state, err := h.gameService.GetState(ctx, gameID, userID)
	if err != nil {
		h.log.Warn().
			Err(err).
			Str("game_id", gameID).
			Str("user_id", userID).
			Msg("websocket initial view denied")
		_ = writeJSON(ctx, wsConn, service.ErrorResponse("", err))
		return
	}

	room := h.hub.GetRoom(gameID)
	room.Add(userID, conn)
	defer func() {
		room.Remove(conn)
		h.broadcastPresence(gameID, userID, room)
		conn.Close()
		h.hub.RemoveRoomIfEmpty(gameID)
	}()

	go conn.WriteLoop(ctx)
	h.broadcastState(room, state, "", nil)

	for {
		msgType, data, err := wsConn.Read(ctx)
		if err != nil {
			h.log.Debug().
				Err(err).
				Str("game_id", gameID).
				Str("user_id", userID).
				Msg("websocket read loop stopped")
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

		newState, events, err := h.gameService.ApplyCommand(ctx, gameID, userID, cmd)
		if err != nil {
			conn.SendJSON(service.ErrorResponse(req.ID, err))
			continue
		}

		h.broadcastState(room, newState, req.ID, events)
	}
}

func (h *Handler) broadcastPresence(gameID string, userID string, room *Room) {
	if room.Size() == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	state, err := h.gameService.GetState(ctx, gameID, userID)
	if err != nil {
		h.log.Debug().
			Err(err).
			Str("game_id", gameID).
			Str("user_id", userID).
			Msg("failed to load state for websocket presence broadcast")
		return
	}

	h.broadcastState(room, state, "", nil)
}

func (h *Handler) broadcastState(room *Room, state domain.GameState, reqID string, events []domain.Event) {
	clients := room.Clients()
	connected := room.ConnectedUserIDs()

	for _, client := range clients {
		view := domain.BuildGameView(state, domain.PlayerID(client.UserID))
		applyConnectedState(&view, connected)
		client.Conn.SendJSON(service.NewUpdateResponse(reqID, view, events))
	}
}

func applyConnectedState(view *domain.GameView, connected map[string]bool) {
	view.Me.Connected = connected[view.Me.UserID.String()]

	for i := range view.Players {
		view.Players[i].Connected = connected[view.Players[i].UserID.String()]
	}
}

func writeJSON(ctx context.Context, wsConn *cws.Conn, v any) error {
	payload, err := json.Marshal(v)
	if err != nil {
		return err
	}

	writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return wsConn.Write(writeCtx, cws.MessageText, payload)
}

func (h *Handler) userIDFromRequest(r *http.Request) (string, error) {
	// 1) Bearer header
	if ah := strings.TrimSpace(r.Header.Get("Authorization")); ah != "" {
		parts := strings.SplitN(ah, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			return "", errors.New("bad authorization header")
		}

		claims, err := h.tokenManager.ParseAccessToken(parts[1])
		if err != nil {
			return "", err
		}
		return claims.UserID, nil
	}

	// 2) Query token
	if token := strings.TrimSpace(r.URL.Query().Get("token")); token != "" {
		claims, err := h.tokenManager.ParseAccessToken(token)
		if err != nil {
			return "", err
		}
		return claims.UserID, nil
	}

	return "", errors.New("missing authorization")
}
