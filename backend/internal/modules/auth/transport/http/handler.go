package http

import (
	"encoding/json"
	"errors"
	"fox/internal/modules/auth/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service      *service.Service
	tokenManager *service.TokenManager
}

func NewHandler(service *service.Service, tokenManager *service.TokenManager) http.Handler {
	h := &Handler{
		service:      service,
		tokenManager: tokenManager,
	}

	r := chi.NewRouter()

	r.Post("/guest", h.Guest)
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Post("/refresh", h.Refresh)

	r.With(AuthMiddleware(h.tokenManager)).Get("/me", h.Me)

	return r
}

func (h *Handler) Guest(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.CreateGuest(r.Context())
	if err != nil {
		http.Error(w, "failed to create guest", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func writeJSON(w http.ResponseWriter, status int, playload any) {
	w.Header().Set("content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(playload)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "username, email and password are required", http.StatusBadRequest)
		return
	}

	result, err := h.service.Register(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrorEmailAlreadyUsed) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "failed to register user", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	result, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrorInvalidCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		http.Error(w, "failed to login user", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaims(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		http.Error(w, "failed to get user", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		http.Error(w, "refresh_token is required", http.StatusBadRequest)
		return
	}

	result, err := h.service.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, service.ErrorInvalidRefreshToken) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if errors.Is(err, service.ErrorRefreshTokenExpired) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		http.Error(w, "failed to refresh access token", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, result)
}
