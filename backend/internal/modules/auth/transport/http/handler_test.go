package http

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	stdhttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"fox/internal/modules/auth/models"
	"fox/internal/modules/auth/repo/postgres"
	"fox/internal/modules/auth/service"
)

type fakeUserRepo struct {
	usersByID    map[string]*models.User
	usersByEmail map[string]*models.User
	nextID       int
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{
		usersByID:    map[string]*models.User{},
		usersByEmail: map[string]*models.User{},
	}
}

func (r *fakeUserRepo) CreateUser(ctx context.Context, params postgres.CreateUserParams) (*models.User, error) {
	r.nextID++

	user := &models.User{
		ID:           fmt.Sprintf("user-%d", r.nextID),
		Username:     params.Username,
		Email:        params.Email,
		PasswordHash: params.PasswordHash,
		IsGuest:      params.IsGuest,
		Role:         params.Role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	r.usersByID[user.ID] = user
	if user.Email != nil {
		r.usersByEmail[*user.Email] = user
	}

	return user, nil
}

func (r *fakeUserRepo) UpdateUser(ctx context.Context, params postgres.UpdateUserParams) error {
	return nil
}

func (r *fakeUserRepo) DeleteUserByID(ctx context.Context, id string) error {
	return nil
}

func (r *fakeUserRepo) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	user, ok := r.usersByID[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return user, nil
}

func (r *fakeUserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, ok := r.usersByEmail[email]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return user, nil
}

type fakeRefreshTokenRepo struct {
	tokens map[string]*models.RefreshToken
	nextID int
}

func newFakeRefreshTokenRepo() *fakeRefreshTokenRepo {
	return &fakeRefreshTokenRepo{
		tokens: map[string]*models.RefreshToken{},
	}
}

func (r *fakeRefreshTokenRepo) CreateRefreshToken(
	ctx context.Context,
	params postgres.CreateRefreshTokenParams,
) (*models.RefreshToken, error) {
	if _, exists := r.tokens[params.Token]; exists {
		return nil, errors.New("refresh token already exists")
	}

	r.nextID++

	token := &models.RefreshToken{
		ID:        fmt.Sprintf("refresh-token-%d", r.nextID),
		UserID:    params.UserID,
		Token:     params.Token,
		ExpiresAt: params.ExpiresAt,
		CreatedAt: time.Now(),
		RevokedAt: nil,
	}

	r.tokens[token.Token] = token

	return token, nil
}

func (r *fakeRefreshTokenRepo) GetRefreshTokenByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	refreshToken, ok := r.tokens[token]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return refreshToken, nil
}

func (r *fakeRefreshTokenRepo) RevokeRefreshTokenByToken(ctx context.Context, token string) error {
	refreshToken, ok := r.tokens[token]
	if !ok {
		return sql.ErrNoRows
	}

	now := time.Now()
	refreshToken.RevokedAt = &now

	return nil
}

func newTestHTTPHandler() (stdhttp.Handler, *service.Service, *service.TokenManager, *fakeUserRepo, *fakeRefreshTokenRepo) {
	userRepo := newFakeUserRepo()
	refreshRepo := newFakeRefreshTokenRepo()
	tm := service.NewTokenManager("test-secret")
	svc := service.NewService(userRepo, refreshRepo, tm)
	handler := NewHandler(svc, tm)

	return handler, svc, tm, userRepo, refreshRepo
}

func performRequest(handler stdhttp.Handler, method, path, body string, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	return rec
}

func TestGuestSuccess(t *testing.T) {
	handler, _, _, _, refreshRepo := newTestHTTPHandler()

	rec := performRequest(handler, stdhttp.MethodPost, "/guest", "", "")

	if rec.Code != stdhttp.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var result service.AuthResult
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if result.User == nil {
		t.Fatal("expected user")
	}

	if !result.User.IsGuest {
		t.Fatal("expected guest user")
	}

	if result.AccessToken == "" {
		t.Fatal("expected access token")
	}

	if result.RefreshToken == "" {
		t.Fatal("expected refresh token")
	}

	if _, ok := refreshRepo.tokens[result.RefreshToken]; !ok {
		t.Fatal("expected refresh token to be saved")
	}
}

func TestRegisterSuccess(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	body := `{
		"username": "alice",
		"email": "alice@example.com",
		"password": "password123"
	}`

	rec := performRequest(handler, stdhttp.MethodPost, "/register", body, "")

	if rec.Code != stdhttp.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var result service.AuthResult
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if result.User == nil {
		t.Fatal("expected user")
	}

	if result.User.Username != "alice" {
		t.Fatalf("expected username alice, got %q", result.User.Username)
	}

	if result.User.Email == nil || *result.User.Email != "alice@example.com" {
		t.Fatal("expected email alice@example.com")
	}

	if result.User.IsGuest {
		t.Fatal("registered user must not be guest")
	}

	if result.AccessToken == "" {
		t.Fatal("expected access token")
	}

	if result.RefreshToken == "" {
		t.Fatal("expected refresh token")
	}
}

func TestRegisterInvalidJSON(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	rec := performRequest(handler, stdhttp.MethodPost, "/register", `{bad json`, "")

	if rec.Code != stdhttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "invalid request body") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestRegisterMissingUsername(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	body := `{
		"email": "alice@example.com",
		"password": "password123"
	}`

	rec := performRequest(handler, stdhttp.MethodPost, "/register", body, "")

	if rec.Code != stdhttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestRegisterMissingEmail(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	body := `{
		"username": "alice",
		"password": "password123"
	}`

	rec := performRequest(handler, stdhttp.MethodPost, "/register", body, "")

	if rec.Code != stdhttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestRegisterMissingPassword(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	body := `{
		"username": "alice",
		"email": "alice@example.com"
	}`

	rec := performRequest(handler, stdhttp.MethodPost, "/register", body, "")

	if rec.Code != stdhttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	body := `{
		"username": "alice",
		"email": "alice@example.com",
		"password": "password123"
	}`

	first := performRequest(handler, stdhttp.MethodPost, "/register", body, "")
	if first.Code != stdhttp.StatusCreated {
		t.Fatalf("expected first register status 201, got %d", first.Code)
	}

	second := performRequest(handler, stdhttp.MethodPost, "/register", body, "")

	if second.Code != stdhttp.StatusConflict {
		t.Fatalf("expected status 409, got %d: %s", second.Code, second.Body.String())
	}
}

func TestLoginSuccess(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	registerBody := `{
		"username": "bob",
		"email": "bob@example.com",
		"password": "password123"
	}`

	registerRec := performRequest(handler, stdhttp.MethodPost, "/register", registerBody, "")
	if registerRec.Code != stdhttp.StatusCreated {
		t.Fatalf("register failed: %d %s", registerRec.Code, registerRec.Body.String())
	}

	loginBody := `{
		"email": "bob@example.com",
		"password": "password123"
	}`

	loginRec := performRequest(handler, stdhttp.MethodPost, "/login", loginBody, "")

	if loginRec.Code != stdhttp.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", loginRec.Code, loginRec.Body.String())
	}

	var result service.AuthResult
	if err := json.NewDecoder(loginRec.Body).Decode(&result); err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if result.User == nil {
		t.Fatal("expected user")
	}

	if result.AccessToken == "" {
		t.Fatal("expected access token")
	}

	if result.RefreshToken == "" {
		t.Fatal("expected refresh token")
	}
}

func TestLoginInvalidJSON(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	rec := performRequest(handler, stdhttp.MethodPost, "/login", `{bad json`, "")

	if rec.Code != stdhttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestLoginMissingEmail(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	body := `{"password":"password123"}`

	rec := performRequest(handler, stdhttp.MethodPost, "/login", body, "")

	if rec.Code != stdhttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestLoginMissingPassword(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	body := `{"email":"bob@example.com"}`

	rec := performRequest(handler, stdhttp.MethodPost, "/login", body, "")

	if rec.Code != stdhttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	registerBody := `{
		"username": "bob",
		"email": "bob@example.com",
		"password": "correct-password"
	}`

	performRequest(handler, stdhttp.MethodPost, "/register", registerBody, "")

	loginBody := `{
		"email": "bob@example.com",
		"password": "wrong-password"
	}`

	rec := performRequest(handler, stdhttp.MethodPost, "/login", loginBody, "")

	if rec.Code != stdhttp.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

func TestLoginUnknownEmail(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	body := `{
		"email": "missing@example.com",
		"password": "password123"
	}`

	rec := performRequest(handler, stdhttp.MethodPost, "/login", body, "")

	if rec.Code != stdhttp.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

func TestRefreshSuccess(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	registerBody := `{
		"username": "charlie",
		"email": "charlie@example.com",
		"password": "password123"
	}`

	registerRec := performRequest(handler, stdhttp.MethodPost, "/register", registerBody, "")
	if registerRec.Code != stdhttp.StatusCreated {
		t.Fatalf("register failed: %d %s", registerRec.Code, registerRec.Body.String())
	}

	var authResult service.AuthResult
	if err := json.NewDecoder(registerRec.Body).Decode(&authResult); err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	refreshBody := fmt.Sprintf(`{"refresh_token":"%s"}`, authResult.RefreshToken)

	refreshRec := performRequest(handler, stdhttp.MethodPost, "/refresh", refreshBody, "")

	if refreshRec.Code != stdhttp.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", refreshRec.Code, refreshRec.Body.String())
	}

	var refreshResult service.RefreshResult
	if err := json.NewDecoder(refreshRec.Body).Decode(&refreshResult); err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if refreshResult.AccessToken == "" {
		t.Fatal("expected access token")
	}

	if refreshResult.RefreshToken == "" {
		t.Fatal("expected refresh token")
	}

	if refreshResult.RefreshToken == authResult.RefreshToken {
		t.Fatal("expected rotated refresh token")
	}
}

func TestRefreshInvalidJSON(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	rec := performRequest(handler, stdhttp.MethodPost, "/refresh", `{bad json`, "")

	if rec.Code != stdhttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestRefreshMissingToken(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	rec := performRequest(handler, stdhttp.MethodPost, "/refresh", `{}`, "")

	if rec.Code != stdhttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestRefreshInvalidToken(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	rec := performRequest(handler, stdhttp.MethodPost, "/refresh", `{"refresh_token":"bad-token"}`, "")

	if rec.Code != stdhttp.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

func TestMeSuccess(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	registerBody := `{
		"username": "dave",
		"email": "dave@example.com",
		"password": "password123"
	}`

	registerRec := performRequest(handler, stdhttp.MethodPost, "/register", registerBody, "")
	if registerRec.Code != stdhttp.StatusCreated {
		t.Fatalf("register failed: %d %s", registerRec.Code, registerRec.Body.String())
	}

	var authResult service.AuthResult
	if err := json.NewDecoder(registerRec.Body).Decode(&authResult); err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	meRec := performRequest(handler, stdhttp.MethodGet, "/me", "", authResult.AccessToken)

	if meRec.Code != stdhttp.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", meRec.Code, meRec.Body.String())
	}

	var user models.User
	if err := json.NewDecoder(meRec.Body).Decode(&user); err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if user.ID != authResult.User.ID {
		t.Fatalf("expected user ID %q, got %q", authResult.User.ID, user.ID)
	}

	if user.Username != "dave" {
		t.Fatalf("expected username dave, got %q", user.Username)
	}
}

func TestMeMissingAuthorizationHeader(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	rec := performRequest(handler, stdhttp.MethodGet, "/me", "", "")

	if rec.Code != stdhttp.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "missing authorization header") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestMeInvalidAuthorizationHeader(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	req := httptest.NewRequest(stdhttp.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Basic abc")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != stdhttp.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

func TestMeInvalidAccessToken(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	rec := performRequest(handler, stdhttp.MethodGet, "/me", "", "invalid-token")

	if rec.Code != stdhttp.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "invalid access token") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestMeTokenForMissingUser(t *testing.T) {
	handler, _, tm, _, _ := newTestHTTPHandler()

	token, err := tm.GenerateAccessToken(&models.User{
		ID:      "missing-user",
		Role:    "player",
		IsGuest: false,
	})
	if err != nil {
		t.Fatalf("GenerateAccessToken returned error: %v", err)
	}

	rec := performRequest(handler, stdhttp.MethodGet, "/me", "", token)

	if rec.Code != stdhttp.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
}

func TestUnknownRoute(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	rec := performRequest(handler, stdhttp.MethodGet, "/unknown", "", "")

	if rec.Code != stdhttp.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

func TestRegisterMethodNotAllowed(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	rec := performRequest(handler, stdhttp.MethodGet, "/register", "", "")

	if rec.Code != stdhttp.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}

func TestLoginMethodNotAllowed(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	rec := performRequest(handler, stdhttp.MethodGet, "/login", "", "")

	if rec.Code != stdhttp.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}

func TestGuestMethodNotAllowed(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	rec := performRequest(handler, stdhttp.MethodGet, "/guest", "", "")

	if rec.Code != stdhttp.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}

func TestRefreshMethodNotAllowed(t *testing.T) {
	handler, _, _, _, _ := newTestHTTPHandler()

	rec := performRequest(handler, stdhttp.MethodGet, "/refresh", "", "")

	if rec.Code != stdhttp.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}

func TestWriteJSONSetsContentType(t *testing.T) {
	rec := httptest.NewRecorder()

	writeJSON(rec, stdhttp.StatusCreated, map[string]string{
		"status": "ok",
	})

	if rec.Code != stdhttp.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}

	contentType := rec.Header().Get("content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected application/json, got %q", contentType)
	}
}

func TestAuthMiddlewarePassesClaimsToNextHandler(t *testing.T) {
	tm := service.NewTokenManager("secret")

	token, err := tm.GenerateAccessToken(&models.User{
		ID:      "user-1",
		Role:    "player",
		IsGuest: false,
	})
	if err != nil {
		t.Fatalf("GenerateAccessToken returned error: %v", err)
	}

	called := false

	next := stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		called = true

		claims, ok := getClaims(r.Context())
		if !ok {
			t.Fatal("expected claims in context")
		}

		if claims.UserID != "user-1" {
			t.Fatalf("expected user-1, got %q", claims.UserID)
		}

		w.WriteHeader(stdhttp.StatusNoContent)
	})

	handler := AuthMiddleware(tm)(next)

	req := httptest.NewRequest(stdhttp.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !called {
		t.Fatal("expected next handler to be called")
	}

	if rec.Code != stdhttp.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
}

func TestAuthMiddlewareRejectsMissingHeader(t *testing.T) {
	tm := service.NewTokenManager("secret")

	next := stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		t.Fatal("next handler must not be called")
	})

	handler := AuthMiddleware(tm)(next)

	req := httptest.NewRequest(stdhttp.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != stdhttp.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

func TestAuthMiddlewareRejectsInvalidPrefix(t *testing.T) {
	tm := service.NewTokenManager("secret")

	next := stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		t.Fatal("next handler must not be called")
	})

	handler := AuthMiddleware(tm)(next)

	req := httptest.NewRequest(stdhttp.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Token abc")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != stdhttp.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

func TestAuthMiddlewareRejectsEmptyBearerToken(t *testing.T) {
	tm := service.NewTokenManager("secret")

	next := stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		t.Fatal("next handler must not be called")
	})

	handler := AuthMiddleware(tm)(next)

	req := httptest.NewRequest(stdhttp.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer ")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != stdhttp.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

func TestAuthMiddlewareRejectsMalformedToken(t *testing.T) {
	tm := service.NewTokenManager("secret")

	next := stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		t.Fatal("next handler must not be called")
	})

	handler := AuthMiddleware(tm)(next)

	req := httptest.NewRequest(stdhttp.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer bad-token")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != stdhttp.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

func TestAuthMiddlewareRejectsTokenSignedWithAnotherSecret(t *testing.T) {
	tm := service.NewTokenManager("secret")
	otherTM := service.NewTokenManager("another-secret")

	token, err := otherTM.GenerateAccessToken(&models.User{
		ID:      "user-1",
		Role:    "player",
		IsGuest: false,
	})
	if err != nil {
		t.Fatalf("GenerateAccessToken returned error: %v", err)
	}

	next := stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		t.Fatal("next handler must not be called")
	})

	handler := AuthMiddleware(tm)(next)

	req := httptest.NewRequest(stdhttp.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != stdhttp.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}
