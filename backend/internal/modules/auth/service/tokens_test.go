package service

import (
	"testing"

	"fox/internal/modules/auth/models"
)

func TestTokenManagerGenerateAndParseAccessToken(t *testing.T) {
	tm := NewTokenManager("test-secret")

	user := &models.User{
		ID:      "user-1",
		Role:    "player",
		IsGuest: false,
	}

	token, err := tm.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("GenerateAccessToken returned error: %v", err)
	}

	if token == "" {
		t.Fatal("access token must not be empty")
	}

	claims, err := tm.ParseAccessToken(token)
	if err != nil {
		t.Fatalf("ParseAccessToken returned error: %v", err)
	}

	if claims.UserID != user.ID {
		t.Fatalf("expected user id %q, got %q", user.ID, claims.UserID)
	}

	if claims.Role != user.Role {
		t.Fatalf("expected role %q, got %q", user.Role, claims.Role)
	}

	if claims.IsGuest != user.IsGuest {
		t.Fatalf("expected is_guest %v, got %v", user.IsGuest, claims.IsGuest)
	}
}

func TestTokenManagerRejectsInvalidAccessToken(t *testing.T) {
	tm := NewTokenManager("test-secret")

	_, err := tm.ParseAccessToken("invalid-token")
	if err == nil {
		t.Fatal("ParseAccessToken should return error for invalid token")
	}
}

func TestTokenManagerGenerateRefreshToken(t *testing.T) {
	tm := NewTokenManager("test-secret")

	token1, err := tm.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken returned error: %v", err)
	}

	token2, err := tm.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken returned error: %v", err)
	}

	if token1 == "" || token2 == "" {
		t.Fatal("refresh token must not be empty")
	}

	if token1 == token2 {
		t.Fatal("refresh tokens must be unique")
	}
}
