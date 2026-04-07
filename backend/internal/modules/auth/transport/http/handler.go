package http

import (
	"encoding/json"
	"errors"
	"fox/internal/modules/auth/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) http.Handler {
	h := &Handler{
		service: service,
	}

	r := chi.NewRouter()

	r.Post("/guest", h.Guest)
	r.Post("/register", h.Register)

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
			w.Write([]byte(err.Error()))
			return
		}
		http.Error(w, "failed to register user", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}
