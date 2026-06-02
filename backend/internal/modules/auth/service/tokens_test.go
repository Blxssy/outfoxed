package service

import (
	"encoding/hex"
	"strings"
	"testing"
	"time"

	"fox/internal/modules/auth/models"

	"github.com/dgrijalva/jwt-go"
)

func newTestUser() *models.User {
	return &models.User{
		ID:      "user-123",
		Role:    "player",
		IsGuest: false,
	}
}

func TestNewTokenManager(t *testing.T) {
	tm := NewTokenManager("secret")

	if tm == nil {
		t.Fatal("expected token manager")
	}

	if string(tm.secret) != "secret" {
		t.Fatalf("expected secret %q, got %q", "secret", string(tm.secret))
	}

	if tm.accessTTL != 15*time.Minute {
		t.Fatalf("expected accessTTL 15 minutes, got %v", tm.accessTTL)
	}
}

func TestGenerateRefreshToken_ReturnsNonEmptyToken(t *testing.T) {
	tm := NewTokenManager("secret")

	token, err := tm.GenerateRefreshToken()

	if err != nil {
		t.Fatalf("GenerateRefreshToken returned error: %v", err)
	}

	if token == "" {
		t.Fatal("refresh token must not be empty")
	}
}

func TestGenerateRefreshToken_HasExpectedLength(t *testing.T) {
	tm := NewTokenManager("secret")

	token, err := tm.GenerateRefreshToken()

	if err != nil {
		t.Fatalf("GenerateRefreshToken returned error: %v", err)
	}

	if len(token) != 64 {
		t.Fatalf("expected token length 64, got %d", len(token))
	}
}

func TestGenerateRefreshToken_IsValidHex(t *testing.T) {
	tm := NewTokenManager("secret")

	token, err := tm.GenerateRefreshToken()

	if err != nil {
		t.Fatalf("GenerateRefreshToken returned error: %v", err)
	}

	decoded, err := hex.DecodeString(token)
	if err != nil {
		t.Fatalf("expected valid hex token, got error: %v", err)
	}

	if len(decoded) != 32 {
		t.Fatalf("expected 32 decoded bytes, got %d", len(decoded))
	}
}

func TestGenerateRefreshToken_GeneratesDifferentTokens(t *testing.T) {
	tm := NewTokenManager("secret")

	token1, err := tm.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("first GenerateRefreshToken returned error: %v", err)
	}

	token2, err := tm.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("second GenerateRefreshToken returned error: %v", err)
	}

	if token1 == token2 {
		t.Fatal("expected different refresh tokens")
	}
}

func TestGenerateRefreshToken_GeneratesManyUniqueTokens(t *testing.T) {
	tm := NewTokenManager("secret")

	tokens := make(map[string]bool)

	for i := 0; i < 100; i++ {
		token, err := tm.GenerateRefreshToken()
		if err != nil {
			t.Fatalf("GenerateRefreshToken returned error: %v", err)
		}

		if tokens[token] {
			t.Fatalf("duplicate refresh token generated: %s", token)
		}

		tokens[token] = true
	}
}

func TestGenerateAccessToken_ReturnsNonEmptyToken(t *testing.T) {
	tm := NewTokenManager("secret")

	token, err := tm.GenerateAccessToken(newTestUser())

	if err != nil {
		t.Fatalf("GenerateAccessToken returned error: %v", err)
	}

	if token == "" {
		t.Fatal("access token must not be empty")
	}
}

func TestGenerateAccessToken_HasThreeJWTParts(t *testing.T) {
	tm := NewTokenManager("secret")

	token, err := tm.GenerateAccessToken(newTestUser())

	if err != nil {
		t.Fatalf("GenerateAccessToken returned error: %v", err)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("expected JWT to have 3 parts, got %d", len(parts))
	}
}

func TestGenerateAccessToken_CanBeParsed(t *testing.T) {
	tm := NewTokenManager("secret")
	user := newTestUser()

	token, err := tm.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("GenerateAccessToken returned error: %v", err)
	}

	claims, err := tm.ParseAccessToken(token)
	if err != nil {
		t.Fatalf("ParseAccessToken returned error: %v", err)
	}

	if claims.UserID != user.ID {
		t.Fatalf("expected user ID %q, got %q", user.ID, claims.UserID)
	}

	if claims.Role != user.Role {
		t.Fatalf("expected role %q, got %q", user.Role, claims.Role)
	}

	if claims.IsGuest != user.IsGuest {
		t.Fatalf("expected isGuest %v, got %v", user.IsGuest, claims.IsGuest)
	}
}

func TestGenerateAccessToken_ForGuestUser(t *testing.T) {
	tm := NewTokenManager("secret")

	user := &models.User{
		ID:      "guest-1",
		Role:    "player",
		IsGuest: true,
	}

	token, err := tm.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("GenerateAccessToken returned error: %v", err)
	}

	claims, err := tm.ParseAccessToken(token)
	if err != nil {
		t.Fatalf("ParseAccessToken returned error: %v", err)
	}

	if claims.UserID != "guest-1" {
		t.Fatalf("expected guest ID, got %q", claims.UserID)
	}

	if !claims.IsGuest {
		t.Fatal("expected IsGuest=true")
	}
}

func TestGenerateAccessToken_ForAdminRole(t *testing.T) {
	tm := NewTokenManager("secret")

	user := &models.User{
		ID:      "admin-1",
		Role:    "admin",
		IsGuest: false,
	}

	token, err := tm.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("GenerateAccessToken returned error: %v", err)
	}

	claims, err := tm.ParseAccessToken(token)
	if err != nil {
		t.Fatalf("ParseAccessToken returned error: %v", err)
	}

	if claims.Role != "admin" {
		t.Fatalf("expected admin role, got %q", claims.Role)
	}
}

func TestGenerateAccessToken_UsesHS256(t *testing.T) {
	tm := NewTokenManager("secret")

	tokenStr, err := tm.GenerateAccessToken(newTestUser())
	if err != nil {
		t.Fatalf("GenerateAccessToken returned error: %v", err)
	}

	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		t.Fatalf("ParseUnverified returned error: %v", err)
	}

	if token.Method != jwt.SigningMethodHS256 {
		t.Fatalf("expected HS256 signing method, got %v", token.Method.Alg())
	}
}

func TestGenerateAccessToken_ContainsExpectedClaims(t *testing.T) {
	tm := NewTokenManager("secret")

	user := &models.User{
		ID:      "user-claims",
		Role:    "moderator",
		IsGuest: true,
	}

	tokenStr, err := tm.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("GenerateAccessToken returned error: %v", err)
	}

	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		t.Fatalf("ParseUnverified returned error: %v", err)
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["sub"] != user.ID {
		t.Fatalf("expected sub %q, got %v", user.ID, claims["sub"])
	}

	if claims["role"] != user.Role {
		t.Fatalf("expected role %q, got %v", user.Role, claims["role"])
	}

	if claims["is_guest"] != user.IsGuest {
		t.Fatalf("expected is_guest %v, got %v", user.IsGuest, claims["is_guest"])
	}

	if _, ok := claims["exp"]; !ok {
		t.Fatal("expected exp claim")
	}
}

func TestGenerateAccessToken_ExpirationIsInFuture(t *testing.T) {
	tm := NewTokenManager("secret")

	before := time.Now().Add(14 * time.Minute).Unix()

	tokenStr, err := tm.GenerateAccessToken(newTestUser())
	if err != nil {
		t.Fatalf("GenerateAccessToken returned error: %v", err)
	}

	after := time.Now().Add(16 * time.Minute).Unix()

	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		t.Fatalf("ParseUnverified returned error: %v", err)
	}

	claims := token.Claims.(jwt.MapClaims)

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		t.Fatalf("expected exp to be float64, got %T", claims["exp"])
	}

	exp := int64(expFloat)

	if exp < before {
		t.Fatalf("expected exp >= %d, got %d", before, exp)
	}

	if exp > after {
		t.Fatalf("expected exp <= %d, got %d", after, exp)
	}
}

func TestParseAccessToken_InvalidTokenString(t *testing.T) {
	tm := NewTokenManager("secret")

	_, err := tm.ParseAccessToken("invalid-token")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseAccessToken_EmptyToken(t *testing.T) {
	tm := NewTokenManager("secret")

	_, err := tm.ParseAccessToken("")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseAccessToken_WrongSecret(t *testing.T) {
	tm1 := NewTokenManager("secret-1")
	tm2 := NewTokenManager("secret-2")

	token, err := tm1.GenerateAccessToken(newTestUser())
	if err != nil {
		t.Fatalf("GenerateAccessToken returned error: %v", err)
	}

	_, err = tm2.ParseAccessToken(token)

	if err == nil {
		t.Fatal("expected error when parsing token with wrong secret")
	}
}

func TestParseAccessToken_TamperedToken(t *testing.T) {
	tm := NewTokenManager("secret")

	token, err := tm.GenerateAccessToken(newTestUser())
	if err != nil {
		t.Fatalf("GenerateAccessToken returned error: %v", err)
	}

	tampered := token + "abc"

	_, err = tm.ParseAccessToken(tampered)

	if err == nil {
		t.Fatal("expected error for tampered token")
	}
}

func TestParseAccessToken_ExpiredToken(t *testing.T) {
	tm := NewTokenManager("secret")

	claims := jwt.MapClaims{
		"sub":      "user-1",
		"role":     "player",
		"is_guest": false,
		"exp":      time.Now().Add(-time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(tm.secret)
	if err != nil {
		t.Fatalf("SignedString returned error: %v", err)
	}

	_, err = tm.ParseAccessToken(tokenStr)

	if err == nil {
		t.Fatal("expected error for expired token")
	}
}

func TestParseAccessToken_MissingSubClaim(t *testing.T) {
	tm := NewTokenManager("secret")

	claims := jwt.MapClaims{
		"role":     "player",
		"is_guest": false,
		"exp":      time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(tm.secret)
	if err != nil {
		t.Fatalf("SignedString returned error: %v", err)
	}

	_, err = tm.ParseAccessToken(tokenStr)

	if err == nil {
		t.Fatal("expected error for missing sub claim")
	}
}

func TestParseAccessToken_InvalidSubClaimType(t *testing.T) {
	tm := NewTokenManager("secret")

	claims := jwt.MapClaims{
		"sub":      123,
		"role":     "player",
		"is_guest": false,
		"exp":      time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(tm.secret)
	if err != nil {
		t.Fatalf("SignedString returned error: %v", err)
	}

	_, err = tm.ParseAccessToken(tokenStr)

	if err == nil {
		t.Fatal("expected error for invalid sub claim")
	}
}

func TestParseAccessToken_MissingRoleClaim(t *testing.T) {
	tm := NewTokenManager("secret")

	claims := jwt.MapClaims{
		"sub":      "user-1",
		"is_guest": false,
		"exp":      time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(tm.secret)
	if err != nil {
		t.Fatalf("SignedString returned error: %v", err)
	}

	_, err = tm.ParseAccessToken(tokenStr)

	if err == nil {
		t.Fatal("expected error for missing role claim")
	}
}

func TestParseAccessToken_InvalidRoleClaimType(t *testing.T) {
	tm := NewTokenManager("secret")

	claims := jwt.MapClaims{
		"sub":      "user-1",
		"role":     123,
		"is_guest": false,
		"exp":      time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(tm.secret)
	if err != nil {
		t.Fatalf("SignedString returned error: %v", err)
	}

	_, err = tm.ParseAccessToken(tokenStr)

	if err == nil {
		t.Fatal("expected error for invalid role claim")
	}
}

func TestParseAccessToken_MissingIsGuestClaim(t *testing.T) {
	tm := NewTokenManager("secret")

	claims := jwt.MapClaims{
		"sub":  "user-1",
		"role": "player",
		"exp":  time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(tm.secret)
	if err != nil {
		t.Fatalf("SignedString returned error: %v", err)
	}

	_, err = tm.ParseAccessToken(tokenStr)

	if err == nil {
		t.Fatal("expected error for missing is_guest claim")
	}
}

func TestParseAccessToken_InvalidIsGuestClaimType(t *testing.T) {
	tm := NewTokenManager("secret")

	claims := jwt.MapClaims{
		"sub":      "user-1",
		"role":     "player",
		"is_guest": "false",
		"exp":      time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(tm.secret)
	if err != nil {
		t.Fatalf("SignedString returned error: %v", err)
	}

	_, err = tm.ParseAccessToken(tokenStr)

	if err == nil {
		t.Fatal("expected error for invalid is_guest claim")
	}
}

func TestParseAccessToken_UnexpectedSigningMethod(t *testing.T) {
	tm := NewTokenManager("secret")

	claims := jwt.MapClaims{
		"sub":      "user-1",
		"role":     "player",
		"is_guest": false,
		"exp":      time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)

	tokenStr, err := token.SignedString(tm.secret)
	if err != nil {
		t.Fatalf("SignedString returned error: %v", err)
	}

	_, err = tm.ParseAccessToken(tokenStr)

	if err == nil {
		t.Fatal("expected error for unexpected signing method")
	}
}

func TestParseAccessToken_AcceptsDifferentRoles(t *testing.T) {
	tm := NewTokenManager("secret")

	roles := []string{
		"player",
		"admin",
		"moderator",
		"",
		"custom-role",
	}

	for _, role := range roles {
		t.Run(role, func(t *testing.T) {
			user := &models.User{
				ID:      "user-role-test",
				Role:    role,
				IsGuest: false,
			}

			token, err := tm.GenerateAccessToken(user)
			if err != nil {
				t.Fatalf("GenerateAccessToken returned error: %v", err)
			}

			claims, err := tm.ParseAccessToken(token)
			if err != nil {
				t.Fatalf("ParseAccessToken returned error: %v", err)
			}

			if claims.Role != role {
				t.Fatalf("expected role %q, got %q", role, claims.Role)
			}
		})
	}
}

func TestParseAccessToken_AcceptsDifferentUserIDs(t *testing.T) {
	tm := NewTokenManager("secret")

	userIDs := []string{
		"user-1",
		"123",
		"",
		"uuid-like-id-0001",
		"very-long-user-id-abcdefghijklmnopqrstuvwxyz",
	}

	for _, userID := range userIDs {
		t.Run(userID, func(t *testing.T) {
			user := &models.User{
				ID:      userID,
				Role:    "player",
				IsGuest: false,
			}

			token, err := tm.GenerateAccessToken(user)
			if err != nil {
				t.Fatalf("GenerateAccessToken returned error: %v", err)
			}

			claims, err := tm.ParseAccessToken(token)
			if err != nil {
				t.Fatalf("ParseAccessToken returned error: %v", err)
			}

			if claims.UserID != userID {
				t.Fatalf("expected user ID %q, got %q", userID, claims.UserID)
			}
		})
	}
}

func TestParseAccessToken_AcceptsGuestAndNonGuestUsers(t *testing.T) {
	tm := NewTokenManager("secret")

	tests := []struct {
		name    string
		isGuest bool
	}{
		{
			name:    "guest",
			isGuest: true,
		},
		{
			name:    "not guest",
			isGuest: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &models.User{
				ID:      "user-1",
				Role:    "player",
				IsGuest: tt.isGuest,
			}

			token, err := tm.GenerateAccessToken(user)
			if err != nil {
				t.Fatalf("GenerateAccessToken returned error: %v", err)
			}

			claims, err := tm.ParseAccessToken(token)
			if err != nil {
				t.Fatalf("ParseAccessToken returned error: %v", err)
			}

			if claims.IsGuest != tt.isGuest {
				t.Fatalf("expected IsGuest %v, got %v", tt.isGuest, claims.IsGuest)
			}
		})
	}
}

func TestParseAccessToken_TableDrivenInvalidClaims(t *testing.T) {
	tm := NewTokenManager("secret")

	tests := []struct {
		name   string
		claims jwt.MapClaims
	}{
		{
			name: "nil sub",
			claims: jwt.MapClaims{
				"sub":      nil,
				"role":     "player",
				"is_guest": false,
				"exp":      time.Now().Add(time.Hour).Unix(),
			},
		},
		{
			name: "nil role",
			claims: jwt.MapClaims{
				"sub":      "user-1",
				"role":     nil,
				"is_guest": false,
				"exp":      time.Now().Add(time.Hour).Unix(),
			},
		},
		{
			name: "nil is_guest",
			claims: jwt.MapClaims{
				"sub":      "user-1",
				"role":     "player",
				"is_guest": nil,
				"exp":      time.Now().Add(time.Hour).Unix(),
			},
		},
		{
			name: "array sub",
			claims: jwt.MapClaims{
				"sub":      []string{"user-1"},
				"role":     "player",
				"is_guest": false,
				"exp":      time.Now().Add(time.Hour).Unix(),
			},
		},
		{
			name: "array role",
			claims: jwt.MapClaims{
				"sub":      "user-1",
				"role":     []string{"player"},
				"is_guest": false,
				"exp":      time.Now().Add(time.Hour).Unix(),
			},
		},
		{
			name: "number is_guest",
			claims: jwt.MapClaims{
				"sub":      "user-1",
				"role":     "player",
				"is_guest": 1,
				"exp":      time.Now().Add(time.Hour).Unix(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, tt.claims)

			tokenStr, err := token.SignedString(tm.secret)
			if err != nil {
				t.Fatalf("SignedString returned error: %v", err)
			}

			_, err = tm.ParseAccessToken(tokenStr)

			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestGenerateAccessToken_ProducesDifferentTokensOverTime(t *testing.T) {
	tm := NewTokenManager("secret")

	user := newTestUser()

	token1, err := tm.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("first GenerateAccessToken returned error: %v", err)
	}

	time.Sleep(time.Second)

	token2, err := tm.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("second GenerateAccessToken returned error: %v", err)
	}

	if token1 == token2 {
		t.Fatal("expected different access tokens generated at different times")
	}
}
