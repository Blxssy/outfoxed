package http

import (
	"encoding/json"
	"testing"
)

func TestRegisterRequestUnmarshalJSON(t *testing.T) {
	data := []byte(`{
		"username": "alice",
		"email": "alice@example.com",
		"password": "password123"
	}`)

	var req RegisterRequest

	err := json.Unmarshal(data, &req)
	if err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if req.Username != "alice" {
		t.Fatalf("expected username alice, got %q", req.Username)
	}

	if req.Email != "alice@example.com" {
		t.Fatalf("expected email alice@example.com, got %q", req.Email)
	}

	if req.Password != "password123" {
		t.Fatalf("expected password password123, got %q", req.Password)
	}
}

func TestRegisterRequestMarshalJSON(t *testing.T) {
	req := RegisterRequest{
		Username: "alice",
		Email:    "alice@example.com",
		Password: "password123",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}

	var decoded map[string]string
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if decoded["username"] != req.Username {
		t.Fatalf("expected username %q, got %q", req.Username, decoded["username"])
	}

	if decoded["email"] != req.Email {
		t.Fatalf("expected email %q, got %q", req.Email, decoded["email"])
	}

	if decoded["password"] != req.Password {
		t.Fatalf("expected password %q, got %q", req.Password, decoded["password"])
	}
}

func TestRegisterRequestZeroValue(t *testing.T) {
	var req RegisterRequest

	if req.Username != "" {
		t.Fatalf("expected empty username, got %q", req.Username)
	}

	if req.Email != "" {
		t.Fatalf("expected empty email, got %q", req.Email)
	}

	if req.Password != "" {
		t.Fatalf("expected empty password, got %q", req.Password)
	}
}

func TestRegisterRequestUnmarshalMissingFields(t *testing.T) {
	data := []byte(`{
		"username": "alice"
	}`)

	var req RegisterRequest

	err := json.Unmarshal(data, &req)
	if err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if req.Username != "alice" {
		t.Fatalf("expected username alice, got %q", req.Username)
	}

	if req.Email != "" {
		t.Fatalf("expected empty email, got %q", req.Email)
	}

	if req.Password != "" {
		t.Fatalf("expected empty password, got %q", req.Password)
	}
}

func TestRegisterRequestUnmarshalUnknownFieldsIgnored(t *testing.T) {
	data := []byte(`{
		"username": "alice",
		"email": "alice@example.com",
		"password": "password123",
		"unknown": "value"
	}`)

	var req RegisterRequest

	err := json.Unmarshal(data, &req)
	if err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if req.Username != "alice" {
		t.Fatalf("expected username alice, got %q", req.Username)
	}

	if req.Email != "alice@example.com" {
		t.Fatalf("expected email alice@example.com, got %q", req.Email)
	}

	if req.Password != "password123" {
		t.Fatalf("expected password password123, got %q", req.Password)
	}
}

func TestLoginRequestUnmarshalJSON(t *testing.T) {
	data := []byte(`{
		"email": "bob@example.com",
		"password": "password123"
	}`)

	var req LoginRequest

	err := json.Unmarshal(data, &req)
	if err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if req.Email != "bob@example.com" {
		t.Fatalf("expected email bob@example.com, got %q", req.Email)
	}

	if req.Password != "password123" {
		t.Fatalf("expected password password123, got %q", req.Password)
	}
}

func TestLoginRequestMarshalJSON(t *testing.T) {
	req := LoginRequest{
		Email:    "bob@example.com",
		Password: "password123",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}

	var decoded map[string]string
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if decoded["email"] != req.Email {
		t.Fatalf("expected email %q, got %q", req.Email, decoded["email"])
	}

	if decoded["password"] != req.Password {
		t.Fatalf("expected password %q, got %q", req.Password, decoded["password"])
	}
}

func TestLoginRequestZeroValue(t *testing.T) {
	var req LoginRequest

	if req.Email != "" {
		t.Fatalf("expected empty email, got %q", req.Email)
	}

	if req.Password != "" {
		t.Fatalf("expected empty password, got %q", req.Password)
	}
}

func TestLoginRequestUnmarshalMissingPassword(t *testing.T) {
	data := []byte(`{
		"email": "bob@example.com"
	}`)

	var req LoginRequest

	err := json.Unmarshal(data, &req)
	if err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if req.Email != "bob@example.com" {
		t.Fatalf("expected email bob@example.com, got %q", req.Email)
	}

	if req.Password != "" {
		t.Fatalf("expected empty password, got %q", req.Password)
	}
}

func TestLoginRequestUnmarshalMissingEmail(t *testing.T) {
	data := []byte(`{
		"password": "password123"
	}`)

	var req LoginRequest

	err := json.Unmarshal(data, &req)
	if err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if req.Email != "" {
		t.Fatalf("expected empty email, got %q", req.Email)
	}

	if req.Password != "password123" {
		t.Fatalf("expected password password123, got %q", req.Password)
	}
}

func TestRefreshRequestUnmarshalJSON(t *testing.T) {
	data := []byte(`{
		"refresh_token": "refresh-token-value"
	}`)

	var req RefreshRequest

	err := json.Unmarshal(data, &req)
	if err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if req.RefreshToken != "refresh-token-value" {
		t.Fatalf("expected refresh-token-value, got %q", req.RefreshToken)
	}
}

func TestRefreshRequestMarshalJSON(t *testing.T) {
	req := RefreshRequest{
		RefreshToken: "refresh-token-value",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}

	var decoded map[string]string
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if decoded["refresh_token"] != req.RefreshToken {
		t.Fatalf("expected refresh_token %q, got %q", req.RefreshToken, decoded["refresh_token"])
	}
}

func TestRefreshRequestZeroValue(t *testing.T) {
	var req RefreshRequest

	if req.RefreshToken != "" {
		t.Fatalf("expected empty refresh token, got %q", req.RefreshToken)
	}
}

func TestRefreshRequestUnmarshalMissingToken(t *testing.T) {
	data := []byte(`{}`)

	var req RefreshRequest

	err := json.Unmarshal(data, &req)
	if err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if req.RefreshToken != "" {
		t.Fatalf("expected empty refresh token, got %q", req.RefreshToken)
	}
}

func TestRequestDTOsInvalidJSON(t *testing.T) {
	tests := []struct {
		name   string
		target any
	}{
		{
			name:   "register request",
			target: &RegisterRequest{},
		},
		{
			name:   "login request",
			target: &LoginRequest{},
		},
		{
			name:   "refresh request",
			target: &RefreshRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := json.Unmarshal([]byte(`{bad json`), tt.target)

			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestRequestDTOsIgnoreUnknownFields(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "register request",
			data: []byte(`{
				"username": "alice",
				"email": "alice@example.com",
				"password": "password123",
				"unknown": "value"
			}`),
		},
		{
			name: "login request",
			data: []byte(`{
				"email": "bob@example.com",
				"password": "password123",
				"unknown": "value"
			}`),
		},
		{
			name: "refresh request",
			data: []byte(`{
				"refresh_token": "refresh-token",
				"unknown": "value"
			}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var decoded map[string]any

			err := json.Unmarshal(tt.data, &decoded)
			if err != nil {
				t.Fatalf("Unmarshal returned error: %v", err)
			}
		})
	}
}
