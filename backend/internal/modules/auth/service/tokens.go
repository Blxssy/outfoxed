package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"fox/internal/modules/auth/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenManager struct {
	secret    []byte
	accessTTL time.Duration
}

type TokenClaims struct {
	UserID  string
	Role    string
	IsGuest bool
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

func (tm *TokenManager) ParseAccessToken(tokenStr string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return tm.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claimsMap, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	userID, ok := claimsMap["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid sub claim")
	}

	role, ok := claimsMap["role"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid role claim")
	}

	isGuest, ok := claimsMap["is_guest"].(bool)
	if !ok {
		return nil, fmt.Errorf("invalid is_guest claim")
	}

	return &TokenClaims{
		UserID:  userID,
		Role:    role,
		IsGuest: isGuest,
	}, nil
}
