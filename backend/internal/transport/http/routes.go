package httptransport

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Deps — зависимости роутера.
type Deps struct {
	// HTTP handlers (обычные REST)
	Auth http.Handler
	Game http.Handler

	// WS handler
	GameWS http.Handler

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// NewRouter собирает chi.Router со всеми путями и middleware.
func NewRouter(d Deps) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	// Healthcheck
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	//WebSocket routes
	r.Route("/ws", func(r chi.Router) {
		r.Get("/games/{id}", d.GameWS.ServeHTTP)
	})

	// API v1
	r.Route("/api/v1", func(r chi.Router) {
		// Таймаут только для обычного HTTP API. Для websocket его нельзя вешать на общий роутер.
		r.Use(middleware.Timeout(15 * time.Second))

		// Монтируем auth
		if d.Auth != nil {
			r.Mount("/auth", d.Auth)
		}

		if d.Game != nil {
			r.Mount("/games", d.Game)
		}
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})

	return r
}
