package http

import (
	"encoding/json"
	"errors"
	authservice "fox/internal/modules/auth/service"
	authhttp "fox/internal/modules/auth/transport/http"
	gameservice "fox/internal/modules/game/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *gameservice.Service
}

func NewHandler(service *gameservice.Service, tokenManager *authservice.TokenManager) http.Handler {
	h := &Handler{service: service}

	r := chi.NewRouter()
	r.Use(authhttp.AuthMiddleware(tokenManager))

	r.Post("/", h.CreateGame)
	r.Get("/{id}", h.GetLobby)
	r.Post("/{id}/join", h.JoinGame)
	r.Post("/{id}/start", h.StartGame)
	r.Get("/{id}/state", h.GetState)

	return r
}

func (h *Handler) CreateGame(w http.ResponseWriter, r *http.Request) {
	var req CreateGameRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}

	userID, ok := userIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	result, err := h.service.CreateGame(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to create game", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func (h *Handler) JoinGame(w http.ResponseWriter, r *http.Request) {
	var req JoinGameRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}

	userID, ok := userIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	gameID := chi.URLParam(r, "id")
	result, err := h.service.JoinGame(r.Context(), gameID, userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) GetLobby(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	gameID := chi.URLParam(r, "id")
	result, err := h.service.GetLobby(r.Context(), gameID, userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) StartGame(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	gameID := chi.URLParam(r, "id")
	result, err := h.service.StartGame(r.Context(), gameID, userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) GetState(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	gameID := chi.URLParam(r, "id")
	result, err := h.service.GetViewState(r.Context(), gameID, userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func userIDFromContext(r *http.Request) (string, bool) {
	claims, ok := authhttp.GetClaims(r.Context())
	if !ok {
		return "", false
	}
	return claims.UserID, true
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, gameservice.ErrForbidden):
		http.Error(w, err.Error(), http.StatusForbidden)
	case errors.Is(err, gameservice.ErrGameNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, gameservice.ErrGameAlreadyStarted),
		errors.Is(err, gameservice.ErrGameFull),
		errors.Is(err, gameservice.ErrNotEnoughPlayers),
		errors.Is(err, gameservice.ErrAlreadyInAnotherGame):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, gameservice.ErrOnlyCreatorCanStart):
		http.Error(w, err.Error(), http.StatusForbidden)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
