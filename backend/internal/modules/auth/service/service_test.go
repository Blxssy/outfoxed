package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"fox/internal/modules/auth/models"
	"fox/internal/modules/auth/repo/postgres"
)

type fakeUserRepo struct {
	usersByID    map[string]*models.User
	usersByEmail map[string]*models.User
	nextID       int
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{
		usersByID:    make(map[string]*models.User),
		usersByEmail: make(map[string]*models.User),
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
	user, ok := r.usersByID[params.ID]
	if !ok {
		return sql.ErrNoRows
	}

	if user.Email != nil {
		delete(r.usersByEmail, *user.Email)
	}

	user.Username = params.Username
	user.Email = params.Email
	user.PasswordHash = params.PasswordHash
	user.IsGuest = params.IsGuest
	user.Role = params.Role
	user.UpdatedAt = time.Now()

	if user.Email != nil {
		r.usersByEmail[*user.Email] = user
	}

	return nil
}

func (r *fakeUserRepo) DeleteUserByID(ctx context.Context, id string) error {
	user, ok := r.usersByID[id]
	if !ok {
		return sql.ErrNoRows
	}

	delete(r.usersByID, id)

	if user.Email != nil {
		delete(r.usersByEmail, *user.Email)
	}

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
		tokens: make(map[string]*models.RefreshToken),
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

func newTestService() (*Service, *fakeUserRepo, *fakeRefreshTokenRepo) {
	userRepo := newFakeUserRepo()
	refreshRepo := newFakeRefreshTokenRepo()
	tokenManager := NewTokenManager("test-secret")

	svc := NewService(userRepo, refreshRepo, tokenManager)

	return svc, userRepo, refreshRepo
}

func TestServiceCreateGuest(t *testing.T) {
	svc, _, refreshRepo := newTestService()

	result, err := svc.CreateGuest(context.Background())
	if err != nil {
		t.Fatalf("CreateGuest returned error: %v", err)
	}

	if result.User == nil {
		t.Fatal("expected user in auth result")
	}

	if !result.User.IsGuest {
		t.Fatal("created user must be guest")
	}

	if result.User.Email != nil {
		t.Fatal("guest user must not have email")
	}

	if result.User.PasswordHash != nil {
		t.Fatal("guest user must not have password hash")
	}

	if result.AccessToken == "" {
		t.Fatal("access token must not be empty")
	}

	if result.RefreshToken == "" {
		t.Fatal("refresh token must not be empty")
	}

	if _, ok := refreshRepo.tokens[result.RefreshToken]; !ok {
		t.Fatal("refresh token must be saved in repository")
	}
}

func TestServiceRegister(t *testing.T) {
	svc, _, refreshRepo := newTestService()

	result, err := svc.Register(
		context.Background(),
		"alice",
		"alice@example.com",
		"password123",
	)
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	if result.User == nil {
		t.Fatal("expected user in auth result")
	}

	if result.User.IsGuest {
		t.Fatal("registered user must not be guest")
	}

	if result.User.Email == nil || *result.User.Email != "alice@example.com" {
		t.Fatal("registered user has invalid email")
	}

	if result.User.PasswordHash == nil {
		t.Fatal("registered user must have password hash")
	}

	if *result.User.PasswordHash == "password123" {
		t.Fatal("password must be stored as hash, not plain text")
	}

	if err := CheckPassword("password123", *result.User.PasswordHash); err != nil {
		t.Fatalf("saved password hash is invalid: %v", err)
	}

	if result.AccessToken == "" {
		t.Fatal("access token must not be empty")
	}

	if result.RefreshToken == "" {
		t.Fatal("refresh token must not be empty")
	}

	if _, ok := refreshRepo.tokens[result.RefreshToken]; !ok {
		t.Fatal("refresh token must be saved in repository")
	}
}

func TestServiceRegisterDuplicateEmail(t *testing.T) {
	svc, _, _ := newTestService()

	_, err := svc.Register(
		context.Background(),
		"alice",
		"alice@example.com",
		"password123",
	)
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	_, err = svc.Register(
		context.Background(),
		"another-alice",
		"alice@example.com",
		"password456",
	)
	if !errors.Is(err, ErrorEmailAlreadyUsed) {
		t.Fatalf("expected ErrorEmailAlreadyUsed, got %v", err)
	}
}

func TestServiceLogin(t *testing.T) {
	svc, userRepo, _ := newTestService()

	email := "bob@example.com"
	password := "password123"

	passwordHash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	_, err = userRepo.CreateUser(context.Background(), postgres.CreateUserParams{
		Username:     "bob",
		Email:        &email,
		PasswordHash: passwordHash,
		IsGuest:      false,
		Role:         "player",
	})
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	result, err := svc.Login(context.Background(), email, password)
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}

	if result.User == nil {
		t.Fatal("expected user in auth result")
	}

	if result.AccessToken == "" {
		t.Fatal("access token must not be empty")
	}

	if result.RefreshToken == "" {
		t.Fatal("refresh token must not be empty")
	}
}

func TestServiceLoginWrongPassword(t *testing.T) {
	svc, userRepo, _ := newTestService()

	email := "bob@example.com"
	passwordHash, err := HashPassword("correct-password")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	_, err = userRepo.CreateUser(context.Background(), postgres.CreateUserParams{
		Username:     "bob",
		Email:        &email,
		PasswordHash: passwordHash,
		IsGuest:      false,
		Role:         "player",
	})
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	_, err = svc.Login(context.Background(), email, "wrong-password")
	if !errors.Is(err, ErrorInvalidCredentials) {
		t.Fatalf("expected ErrorInvalidCredentials, got %v", err)
	}
}

func TestServiceRefreshRotatesToken(t *testing.T) {
	svc, userRepo, refreshRepo := newTestService()

	email := "charlie@example.com"

	user, err := userRepo.CreateUser(context.Background(), postgres.CreateUserParams{
		Username:     "charlie",
		Email:        &email,
		PasswordHash: nil,
		IsGuest:      false,
		Role:         "player",
	})
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	oldRefreshToken := "old-refresh-token"

	_, err = refreshRepo.CreateRefreshToken(context.Background(), postgres.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     oldRefreshToken,
		ExpiresAt: time.Now().Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("CreateRefreshToken returned error: %v", err)
	}

	result, err := svc.Refresh(context.Background(), oldRefreshToken)
	if err != nil {
		t.Fatalf("Refresh returned error: %v", err)
	}

	if result.AccessToken == "" {
		t.Fatal("new access token must not be empty")
	}

	if result.RefreshToken == "" {
		t.Fatal("new refresh token must not be empty")
	}

	if result.RefreshToken == oldRefreshToken {
		t.Fatal("new refresh token must be different from old refresh token")
	}

	oldTokenRecord, err := refreshRepo.GetRefreshTokenByToken(context.Background(), oldRefreshToken)
	if err != nil {
		t.Fatalf("GetRefreshTokenByToken returned error: %v", err)
	}

	if oldTokenRecord.RevokedAt == nil {
		t.Fatal("old refresh token must be revoked")
	}

	if _, ok := refreshRepo.tokens[result.RefreshToken]; !ok {
		t.Fatal("new refresh token must be saved in repository")
	}
}

func TestServiceRefreshExpiredToken(t *testing.T) {
	svc, userRepo, refreshRepo := newTestService()

	email := "dave@example.com"

	user, err := userRepo.CreateUser(context.Background(), postgres.CreateUserParams{
		Username:     "dave",
		Email:        &email,
		PasswordHash: nil,
		IsGuest:      false,
		Role:         "player",
	})
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	expiredToken := "expired-refresh-token"

	_, err = refreshRepo.CreateRefreshToken(context.Background(), postgres.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     expiredToken,
		ExpiresAt: time.Now().Add(-time.Hour),
	})
	if err != nil {
		t.Fatalf("CreateRefreshToken returned error: %v", err)
	}

	_, err = svc.Refresh(context.Background(), expiredToken)
	if !errors.Is(err, ErrorRefreshTokenExpired) {
		t.Fatalf("expected ErrorRefreshTokenExpired, got %v", err)
	}
}

func TestServiceRefreshRevokedToken(t *testing.T) {
	svc, userRepo, refreshRepo := newTestService()

	email := "eve@example.com"

	user, err := userRepo.CreateUser(context.Background(), postgres.CreateUserParams{
		Username:     "eve",
		Email:        &email,
		PasswordHash: nil,
		IsGuest:      false,
		Role:         "player",
	})
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	revokedToken := "revoked-refresh-token"

	_, err = refreshRepo.CreateRefreshToken(context.Background(), postgres.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     revokedToken,
		ExpiresAt: time.Now().Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("CreateRefreshToken returned error: %v", err)
	}

	err = refreshRepo.RevokeRefreshTokenByToken(context.Background(), revokedToken)
	if err != nil {
		t.Fatalf("RevokeRefreshTokenByToken returned error: %v", err)
	}

	_, err = svc.Refresh(context.Background(), revokedToken)
	if !errors.Is(err, ErrorInvalidRefreshToken) {
		t.Fatalf("expected ErrorInvalidRefreshToken, got %v", err)
	}
}

func TestServiceLoginUserNotFound(t *testing.T) {
	svc, _, _ := newTestService()

	_, err := svc.Login(context.Background(), "missing@example.com", "password")

	if !errors.Is(err, ErrorInvalidCredentials) {
		t.Fatalf("expected ErrorInvalidCredentials, got %v", err)
	}
}

func TestServiceLoginGuestUser(t *testing.T) {
	svc, userRepo, _ := newTestService()

	email := "guest@example.com"

	_, err := userRepo.CreateUser(context.Background(), postgres.CreateUserParams{
		Username:     "guest",
		Email:        &email,
		PasswordHash: nil,
		IsGuest:      true,
		Role:         "player",
	})
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	_, err = svc.Login(context.Background(), email, "password")

	if !errors.Is(err, ErrorInvalidCredentials) {
		t.Fatalf("expected ErrorInvalidCredentials, got %v", err)
	}
}

func TestServiceLoginUserWithoutPasswordHash(t *testing.T) {
	svc, userRepo, _ := newTestService()

	email := "nopassword@example.com"

	_, err := userRepo.CreateUser(context.Background(), postgres.CreateUserParams{
		Username:     "nopassword",
		Email:        &email,
		PasswordHash: nil,
		IsGuest:      false,
		Role:         "player",
	})
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	_, err = svc.Login(context.Background(), email, "password")

	if !errors.Is(err, ErrorInvalidCredentials) {
		t.Fatalf("expected ErrorInvalidCredentials, got %v", err)
	}
}

func TestServiceGetUserByIDSuccess(t *testing.T) {
	svc, userRepo, _ := newTestService()

	email := "user@example.com"

	createdUser, err := userRepo.CreateUser(context.Background(), postgres.CreateUserParams{
		Username:     "user",
		Email:        &email,
		PasswordHash: nil,
		IsGuest:      false,
		Role:         "player",
	})
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	user, err := svc.GetUserByID(context.Background(), createdUser.ID)

	if err != nil {
		t.Fatalf("GetUserByID returned error: %v", err)
	}

	if user.ID != createdUser.ID {
		t.Fatalf("expected user ID %q, got %q", createdUser.ID, user.ID)
	}
}

func TestServiceGetUserByIDNotFound(t *testing.T) {
	svc, _, _ := newTestService()

	_, err := svc.GetUserByID(context.Background(), "missing-id")

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestServiceRefreshMissingToken(t *testing.T) {
	svc, _, _ := newTestService()

	_, err := svc.Refresh(context.Background(), "missing-refresh-token")

	if !errors.Is(err, ErrorInvalidRefreshToken) {
		t.Fatalf("expected ErrorInvalidRefreshToken, got %v", err)
	}
}

func TestServiceRefreshTokenUserNotFound(t *testing.T) {
	svc, _, refreshRepo := newTestService()

	token := "token-with-missing-user"

	_, err := refreshRepo.CreateRefreshToken(context.Background(), postgres.CreateRefreshTokenParams{
		UserID:    "missing-user-id",
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("CreateRefreshToken returned error: %v", err)
	}

	_, err = svc.Refresh(context.Background(), token)

	if !errors.Is(err, ErrorInvalidRefreshToken) {
		t.Fatalf("expected ErrorInvalidRefreshToken, got %v", err)
	}
}

func TestServiceRegisterCreatesUserWithCorrectFields(t *testing.T) {
	svc, _, _ := newTestService()

	result, err := svc.Register(
		context.Background(),
		"john",
		"john@example.com",
		"password123",
	)

	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	if result.User.Username != "john" {
		t.Fatalf("expected username john, got %q", result.User.Username)
	}

	if result.User.Role != "player" {
		t.Fatalf("expected role player, got %q", result.User.Role)
	}

	if result.User.IsGuest {
		t.Fatal("registered user must not be guest")
	}
}

func TestServiceCreateGuestCreatesUserWithCorrectFields(t *testing.T) {
	svc, _, _ := newTestService()

	result, err := svc.CreateGuest(context.Background())

	if err != nil {
		t.Fatalf("CreateGuest returned error: %v", err)
	}

	if result.User.Username == "" {
		t.Fatal("guest username must not be empty")
	}

	if result.User.Role != "player" {
		t.Fatalf("expected role player, got %q", result.User.Role)
	}

	if !result.User.IsGuest {
		t.Fatal("guest user must have IsGuest=true")
	}
}

func TestServiceRegisterDifferentUsersGetDifferentRefreshTokens(t *testing.T) {
	svc, _, _ := newTestService()

	first, err := svc.Register(
		context.Background(),
		"first",
		"first@example.com",
		"password123",
	)
	if err != nil {
		t.Fatalf("first Register returned error: %v", err)
	}

	second, err := svc.Register(
		context.Background(),
		"second",
		"second@example.com",
		"password123",
	)
	if err != nil {
		t.Fatalf("second Register returned error: %v", err)
	}

	if first.RefreshToken == second.RefreshToken {
		t.Fatal("different users must receive different refresh tokens")
	}
}

func TestServiceLoginCreatesNewRefreshTokenEachTime(t *testing.T) {
	svc, _, _ := newTestService()

	email := "login@example.com"
	password := "password123"

	_, err := svc.Register(context.Background(), "login-user", email, password)
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	firstLogin, err := svc.Login(context.Background(), email, password)
	if err != nil {
		t.Fatalf("first Login returned error: %v", err)
	}

	secondLogin, err := svc.Login(context.Background(), email, password)
	if err != nil {
		t.Fatalf("second Login returned error: %v", err)
	}

	if firstLogin.RefreshToken == secondLogin.RefreshToken {
		t.Fatal("each login must create a new refresh token")
	}
}
