package postgres

import (
	"context"
	"fox/internal/modules/auth/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type CreateRefreshTokenParams struct {
	UserID    string
	Token     string
	ExpiresAt time.Time
}

type RefreshTokenRepo interface {
	CreateRefreshToken(ctx context.Context, params CreateRefreshTokenParams) (*models.RefreshToken, error)
	GetRefreshTokenByToken(ctx context.Context, token string) (*models.RefreshToken, error)
	RevokeRefreshTokenByToken(ctx context.Context, token string) error
}

type RefreshTokenRepository struct {
	db *sqlx.DB
}

func NewRefreshTokenRepository(db *sqlx.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db: db,
	}
}

func (r *RefreshTokenRepository) CreateRefreshToken(
	ctx context.Context,
	params CreateRefreshTokenParams,
) (*models.RefreshToken, error) {
	query := `
		INSERT INTO refresh_tokens (
			user_id,
			token,
			expires_at
		)
		VALUES (
			$1, $2, $3
		)
		RETURNING
			id,
			user_id,
			token,
			expires_at,
			created_at,
			revoked_at
	`

	var refreshToken models.RefreshToken

	err := r.db.GetContext(
		ctx,
		&refreshToken,
		query,
		params.UserID,
		params.Token,
		params.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (r *RefreshTokenRepository) GetRefreshTokenByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	query := `
		SELECT
			id,
			user_id,
			token,
			expires_at,
			created_at,
			revoked_at
		FROM refresh_tokens
		WHERE token = $1
	`

	var refreshToken models.RefreshToken

	err := r.db.GetContext(ctx, &refreshToken, query, token)
	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (r *RefreshTokenRepository) RevokeRefreshTokenByToken(ctx context.Context, token string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE token = $1 AND revoked_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query, token)
	return err
}
