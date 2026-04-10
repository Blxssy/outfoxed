package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"fox/internal/modules/auth/models"
	"fox/internal/modules/auth/repo/postgres"
	"time"
)

var ErrorEmailAlreadyUsed = errors.New("The email address has already been used")
var ErrorInvalidCredentials = errors.New("invalid email or password")
var ErrorInvalidRefreshToken = errors.New("invalid refresh token")
var ErrorRefreshTokenExpired = errors.New("refresh token expired")

type Service struct {
	userRepo         postgres.UserRepo
	refreshTokenRepo postgres.RefreshTokenRepo
	tokenManager     *TokenManager
}

func NewService(
	userRepo postgres.UserRepo,
	refreshTokenRepo postgres.RefreshTokenRepo,
	tm *TokenManager,
) *Service {
	return &Service{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		tokenManager:     tm,
	}
}

func generateGuestUsername() string {
	return fmt.Sprintf("Guest_%d", time.Now().UnixNano())
}

func (s *Service) CreateGuest(ctx context.Context) (*AuthResult, error) {
	username := generateGuestUsername()

	user, err := s.userRepo.CreateUser(ctx, postgres.CreateUserParams{
		Username:     username,
		Email:        nil,
		PasswordHash: nil,
		IsGuest:      true,
		Role:         "player",
	})
	if err != nil {
		return nil, err
	}

	accessToken, err := s.tokenManager.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	_, err = s.refreshTokenRepo.CreateRefreshToken(ctx, postgres.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	})
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Register(ctx context.Context, username, email, password string) (*AuthResult, error) {
	existingUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, ErrorEmailAlreadyUsed
	}
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.CreateUser(ctx, postgres.CreateUserParams{
		Username:     username,
		Email:        &email,
		PasswordHash: passwordHash,
		IsGuest:      false,
		Role:         "player",
	})
	if err != nil {
		return nil, err
	}

	accessToken, err := s.tokenManager.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	_, err = s.refreshTokenRepo.CreateRefreshToken(ctx, postgres.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	})
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorInvalidCredentials
		}
		return nil, err
	}

	if user.IsGuest || user.PasswordHash == nil {
		return nil, ErrorInvalidCredentials
	}

	if err := CheckPassword(password, *user.PasswordHash); err != nil {
		return nil, ErrorInvalidCredentials
	}

	accessToken, err := s.tokenManager.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*RefreshResult, error) {
	tokenRecord, err := s.refreshTokenRepo.GetRefreshTokenByToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorInvalidRefreshToken
		}
		return nil, err
	}

	if tokenRecord.RevokedAt != nil {
		return nil, ErrorInvalidRefreshToken
	}

	if time.Now().After(tokenRecord.ExpiresAt) {
		return nil, ErrorRefreshTokenExpired
	}

	user, err := s.userRepo.GetUserByID(ctx, tokenRecord.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorInvalidRefreshToken
		}
		return nil, err
	}

	newAccessToken, err := s.tokenManager.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.tokenManager.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	err = s.refreshTokenRepo.RevokeRefreshTokenByToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	_, err = s.refreshTokenRepo.CreateRefreshToken(ctx, postgres.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	})
	if err != nil {
		return nil, err
	}

	return &RefreshResult{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
