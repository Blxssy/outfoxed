package http

import (
	"context"

	"fox/internal/modules/auth/service"
)

type contextKey string

const claimsContextKey contextKey = "auth_claims"

func withClaims(ctx context.Context, claims *service.TokenClaims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}

func getClaims(ctx context.Context) (*service.TokenClaims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*service.TokenClaims)
	return claims, ok
}
