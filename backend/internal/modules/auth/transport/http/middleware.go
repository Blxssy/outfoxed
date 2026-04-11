package http

import (
	"net/http"
	"strings"

	"fox/internal/modules/auth/service"
)

func AuthMiddleware(tm *service.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			const prefix = "Bearer "
			if !strings.HasPrefix(authHeader, prefix) {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, prefix)

			claims, err := tm.ParseAccessToken(token)
			if err != nil {
				http.Error(w, "invalid access token", http.StatusUnauthorized)
				return
			}

			ctx := withClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
