package service

import (
	"context"
	"fmt"
	"fox/internal/modules/auth/repo/postgres"
	"time"
)

type Service struct {
	repo         postgres.UserRepo
	tokenManager *TokenManager
}

func NewService(repo postgres.UserRepo, tm *TokenManager) *Service {
	return &Service{
		repo:         repo,
		tokenManager: tm,
	}
}

func generateGuestUsername() string {
	return fmt.Sprintf("Guest_%d", time.Now().UnixNano())
}

func (s *Service) CreateGuest(ctx context.Context) (*AuthResult, error) {
	username := generateGuestUsername()

	user, err := s.repo.CreateUser(ctx, postgres.CreateUserParams{
		Username:     username,
		Email:        nil,
		PasswordHash: nil,
		IsGuest:      true,
		Role:         "player",
	})
	if err != nil {
		return nil, fmt.Errorf("s.repo.CreateUser: %w", err)
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
