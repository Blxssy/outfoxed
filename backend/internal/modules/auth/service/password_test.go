package service

import (
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword_ReturnsHash(t *testing.T) {
	password := "my-secret-password"

	hash, err := HashPassword(password)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if hash == nil {
		t.Fatal("expected hash pointer, got nil")
	}
	if *hash == "" {
		t.Fatal("expected non-empty hash")
	}
}

func TestHashPassword_HashIsNotPlainPassword(t *testing.T) {
	password := "password123"

	hash, err := HashPassword(password)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if *hash == password {
		t.Fatal("hash must not equal plain password")
	}
}

func TestHashPassword_ReturnsValidBcryptHash(t *testing.T) {
	password := "strong-password"

	hash, err := HashPassword(password)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.HasPrefix(*hash, "$2a$") &&
		!strings.HasPrefix(*hash, "$2b$") &&
		!strings.HasPrefix(*hash, "$2y$") {
		t.Fatalf("expected bcrypt hash prefix, got %q", *hash)
	}
}

func TestHashPassword_CanBeCheckedWithCorrectPassword(t *testing.T) {
	password := "correct-password"

	hash, err := HashPassword(password)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = CheckPassword(password, *hash)

	if err != nil {
		t.Fatalf("expected password check to pass, got %v", err)
	}
}

func TestCheckPassword_FailsWithWrongPassword(t *testing.T) {
	password := "correct-password"
	wrongPassword := "wrong-password"

	hash, err := HashPassword(password)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = CheckPassword(wrongPassword, *hash)

	if err == nil {
		t.Fatal("expected password check to fail")
	}
}

func TestCheckPassword_FailsWithEmptyPassword(t *testing.T) {
	password := "not-empty"

	hash, err := HashPassword(password)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = CheckPassword("", *hash)

	if err == nil {
		t.Fatal("expected empty password to fail")
	}
}

func TestCheckPassword_WorksWithEmptyOriginalPassword(t *testing.T) {
	password := ""

	hash, err := HashPassword(password)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = CheckPassword(password, *hash)

	if err != nil {
		t.Fatalf("expected empty original password to pass, got %v", err)
	}
}

func TestCheckPassword_FailsWithInvalidHash(t *testing.T) {
	err := CheckPassword("password", "invalid-hash")

	if err == nil {
		t.Fatal("expected invalid hash to return error")
	}
}

func TestHashPassword_GeneratesDifferentHashesForSamePassword(t *testing.T) {
	password := "same-password"

	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if *hash1 == *hash2 {
		t.Fatal("expected different hashes because bcrypt uses salt")
	}
}

func TestHashPassword_BcryptCostIsDefaultCost(t *testing.T) {
	password := "cost-password"

	hash, err := HashPassword(password)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	cost, err := bcrypt.Cost([]byte(*hash))
	if err != nil {
		t.Fatalf("expected valid bcrypt cost, got %v", err)
	}

	if cost != bcrypt.DefaultCost {
		t.Fatalf("expected cost %d, got %d", bcrypt.DefaultCost, cost)
	}
}

func TestCheckPassword_TableDriven(t *testing.T) {
	originalPassword := "table-password"

	hash, err := HashPassword(originalPassword)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	tests := []struct {
		name        string
		password    string
		hash        string
		expectError bool
	}{
		{
			name:        "correct password",
			password:    originalPassword,
			hash:        *hash,
			expectError: false,
		},
		{
			name:        "wrong password",
			password:    "wrong",
			hash:        *hash,
			expectError: true,
		},
		{
			name:        "empty password",
			password:    "",
			hash:        *hash,
			expectError: true,
		},
		{
			name:        "empty hash",
			password:    originalPassword,
			hash:        "",
			expectError: true,
		},
		{
			name:        "invalid hash",
			password:    originalPassword,
			hash:        "not-a-bcrypt-hash",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPassword(tt.password, tt.hash)

			if tt.expectError && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tt.expectError && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}
