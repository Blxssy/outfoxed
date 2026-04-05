package postgres

import (
	"context"

	"fox/internal/modules/auth/models"

	"github.com/jmoiron/sqlx"
)

type CreateUserParams struct {
	Username     string
	Email        *string
	PasswordHash *string
	IsGuest      bool
	Role         string
}

type UpdateUserParams struct {
	ID           string
	Username     string
	Email        *string
	PasswordHash *string
	IsGuest      bool
	Role         string
}

type UserRepo interface {
	CreateUser(ctx context.Context, params CreateUserParams) error
	UpdateUser(ctx context.Context, params UpdateUserParams) error
	DeleteUserByID(ctx context.Context, id string) error

	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) CreateUser(ctx context.Context, params CreateUserParams) error {
	query := `
	INSERT INTO users (
		username,
		email,
		password_hash,
		is_guest,
		role
	)
	VALUES (
		$1, $2, $3, $4, $5
	)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		params.Username,
		params.Email,
		params.PasswordHash,
		params.IsGuest,
		params.Role,
	)
	return err
}

func (r *userRepo) UpdateUser(ctx context.Context, params UpdateUserParams) error {
	query := `
	UPDATE users 
	SET 
		username = $1,
		email = $2,
		password_hash = $3,
		is_guest = $4,
		role = $5
	WHERE id = $6
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		params.Username,
		params.Email,
		params.PasswordHash,
		params.IsGuest,
		params.Role,
		params.ID,
	)

	return err
}

func (r *userRepo) DeleteUserByID(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id=$1`

	_, err := r.db.ExecContext(
		ctx,
		query,
		id,
	)
	return err
}

func (r *userRepo) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `
	SELECT 
		id, 
		username, 
		email, 
		password_hash, 
		is_guest, 
		role, 
		created_at, 
		updated_at, 
		last_seen_at
	FROM users 
	WHERE id=$1`
	var user models.User
	err := r.db.GetContext(
		ctx,
		&user,
		query,
		id,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
	SELECT 
		id, 
		username, 
		email, 
		password_hash, 
		is_guest, 
		role, 
		created_at, 
		updated_at, 
		last_seen_at
	FROM users 
	WHERE id=$1`
	var user models.User
	err := r.db.GetContext(
		ctx,
		&user,
		query,
		email,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
