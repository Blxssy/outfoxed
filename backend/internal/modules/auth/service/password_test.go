package service

import "testing"

func TestHashPasswordAndCheckPassword(t *testing.T) {
	password := "strong-password-123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if hash == nil {
		t.Fatal("HashPassword returned nil hash")
	}

	if *hash == password {
		t.Fatal("password hash must not be equal to original password")
	}

	if err := CheckPassword(password, *hash); err != nil {
		t.Fatalf("CheckPassword failed for correct password: %v", err)
	}
}

func TestCheckPasswordWrongPassword(t *testing.T) {
	password := "correct-password"
	wrongPassword := "wrong-password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if err := CheckPassword(wrongPassword, *hash); err == nil {
		t.Fatal("CheckPassword should return error for wrong password")
	}
}
