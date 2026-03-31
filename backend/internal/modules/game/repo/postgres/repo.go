package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"fox/internal/modules/game/repo"
)

type Repo struct {
	db *sql.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *Repo) IsPlayerInGame(ctx context.Context, tx *sql.Tx, gameID string, userID string) (bool, error) {
	var exists bool
	err := tx.QueryRowContext(ctx, `
		select exists(
			select 1 from game_players
			where game_id = $1 and user_id = $2
		)
	`, gameID, userID).Scan(&exists)
	return exists, err
}

func (r *Repo) IsPlayerInGameReadonly(ctx context.Context, gameID string, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		select exists(
			select 1 from game_players
			where game_id = $1 and user_id = $2
		)
	`, gameID, userID).Scan(&exists)
	return exists, err
}

func (r *Repo) GetGameForUpdate(ctx context.Context, tx *sql.Tx, gameID string) (repo.GameRow, error) {
	var row repo.GameRow
	err := tx.QueryRowContext(ctx, `
		select id, status, state_json, version, fox_escape_at, culprit_id
		from games
		where id = $1
		for update
	`, gameID).Scan(&row.ID, &row.Status, &row.StateJSON, &row.Version, &row.FoxEscapeAt, &row.CulpritID)
	if err != nil {
		return repo.GameRow{}, err
	}
	return row, nil
}

func (r *Repo) GetGame(ctx context.Context, gameID string) (repo.GameRow, error) {
	var row repo.GameRow
	err := r.db.QueryRowContext(ctx, `
		select id, status, state_json, version, fox_escape_at, culprit_id
		from games
		where id = $1
	`, gameID).Scan(&row.ID, &row.Status, &row.StateJSON, &row.Version, &row.FoxEscapeAt, &row.CulpritID)
	if err != nil {
		return repo.GameRow{}, err
	}
	return row, nil
}

func (r *Repo) UpdateState(ctx context.Context, tx *sql.Tx, gameID string, newStateJSON []byte, newVersion int) error {
	res, err := tx.ExecContext(ctx, `
		update games
		set state_json = $2,
		    version = $3,
		    updated_at = now()
		where id = $1
	`, gameID, newStateJSON, newVersion)
	if err != nil {
		return err
	}

	aff, _ := res.RowsAffected()
	if aff != 1 {
		return fmt.Errorf("update affected %d rows", aff)
	}
	return nil
}
