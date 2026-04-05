package models

import "time"

type User struct {
	ID           string     `json:"id" db:"id"`
	Username     string     `json:"username" db:"username"`
	Email        *string    `json:"email,omitempty" db:"email"`
	PasswordHash *string    `json:"-" db:"password_hash"`
	IsGuest      bool       `json:"is_guest" db:"is_guest"`
	Role         string     `json:"role" db:"role"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
}
