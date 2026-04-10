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

	initialView, err := h.gameService.GetView(ctx, gameID, userID)
	if err != nil {
		conn.SendJSON(service.NewErrorResponse("", "forbidden", "You do not have access to this game."))
		return
	}

	room := h.hub.GetRoom(gameID)
	room.Add(userID, conn)
	defer func() {
		room.Remove(conn)
		conn.Close()
		h.hub.RemoveRoomIfEmpty(gameID)
	}()

	go conn.WriteLoop(ctx)

	conn.SendJSON(service.NewUpdateResponse("", initialView, nil))

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

		newState, events, err := h.gameService.ApplyCommand(ctx, gameID, userID, cmd)
		if err != nil {
			conn.SendJSON(service.ErrorResponse(req.ID, err))
			continue
		}

		clients := room.Clients()
		for _, client := range clients {
			view := domain.BuildGameView(newState, domain.PlayerID(client.UserID))
			client.Conn.SendJSON(service.NewUpdateResponse(req.ID, view, events))
		}
	}
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
