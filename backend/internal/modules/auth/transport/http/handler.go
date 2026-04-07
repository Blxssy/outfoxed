package http

import (
	"encoding/json"
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
