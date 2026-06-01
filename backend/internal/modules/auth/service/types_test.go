package service

import (
	"encoding/json"
	"testing"

	"fox/internal/modules/auth/models"
)

func TestAuthResultJSON(t *testing.T) {
	result := AuthResult{
		User: &models.User{
			ID:       "user-1",
			Username: "alice",
			Role:     "player",
			IsGuest:  false,
		},
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}

	jsonStr := string(data)

	expectedParts := []string{
		`"user"`,
		`"access_token":"access-token"`,
		`"refresh_token":"refresh-token"`,
	}

	for _, part := range expectedParts {
		if !contains(jsonStr, part) {
			t.Fatalf("expected JSON to contain %s, got %s", part, jsonStr)
		}
	}
}

func TestAuthResultJSONWithNilUser(t *testing.T) {
	result := AuthResult{
		User:         nil,
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if decoded["user"] != nil {
		t.Fatalf("expected user to be nil, got %v", decoded["user"])
	}

	if decoded["access_token"] != "access-token" {
		t.Fatalf("expected access_token, got %v", decoded["access_token"])
	}

	if decoded["refresh_token"] != "refresh-token" {
		t.Fatalf("expected refresh_token, got %v", decoded["refresh_token"])
	}
}

func TestRefreshResultJSON(t *testing.T) {
	result := RefreshResult{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}

	var decoded map[string]string
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if decoded["access_token"] != result.AccessToken {
		t.Fatalf("expected access_token %q, got %q", result.AccessToken, decoded["access_token"])
	}

	if decoded["refresh_token"] != result.RefreshToken {
		t.Fatalf("expected refresh_token %q, got %q", result.RefreshToken, decoded["refresh_token"])
	}
}

func TestRefreshResultUnmarshalJSON(t *testing.T) {
	jsonData := []byte(`{
		"access_token": "access-token",
		"refresh_token": "refresh-token"
	}`)

	var result RefreshResult

	err := json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if result.AccessToken != "access-token" {
		t.Fatalf("expected access-token, got %q", result.AccessToken)
	}

	if result.RefreshToken != "refresh-token" {
		t.Fatalf("expected refresh-token, got %q", result.RefreshToken)
	}
}

func TestAuthResultZeroValue(t *testing.T) {
	var result AuthResult

	if result.User != nil {
		t.Fatal("expected zero value User to be nil")
	}

	if result.AccessToken != "" {
		t.Fatal("expected zero value AccessToken to be empty")
	}

	if result.RefreshToken != "" {
		t.Fatal("expected zero value RefreshToken to be empty")
	}
}

func TestRefreshResultZeroValue(t *testing.T) {
	var result RefreshResult

	if result.AccessToken != "" {
		t.Fatal("expected zero value AccessToken to be empty")
	}

	if result.RefreshToken != "" {
		t.Fatal("expected zero value RefreshToken to be empty")
	}
}

func contains(s, substr string) bool {
	return len(substr) == 0 || len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
