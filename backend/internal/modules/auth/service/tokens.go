package service

import (
	"crypto/rand"
	"encoding/hex"
	"fox/internal/modules/auth/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenManager struct {
	secret    []byte
	accessTTL time.Duration
}

func NewTokenManager(secret string) *TokenManager {
	return &TokenManager{
		secret:    []byte(secret),
		accessTTL: 15 * time.Minute,
	}
}

func (tm *TokenManager) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func (tm *TokenManager) GenerateAccessToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":      user.ID,
		"role":     user.Role,
		"is_guest": user.IsGuest,
		"exp":      time.Now().Add(tm.accessTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(tm.secret)
}
