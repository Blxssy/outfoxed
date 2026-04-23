package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"fox/internal/modules/game/repo"

	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *Repo) CreateGame(
	ctx context.Context,
	tx *sql.Tx,
	createdBy string,
	title string,
	visibility string,
	joinCode *string,
	stateJSON []byte,
) (repo.GameRow, error) {
	var row repo.GameRow

	var joinCodeValue any
	if joinCode != nil {
		joinCodeValue = *joinCode
	}

	err := tx.QueryRowContext(ctx, `
		insert into games (
			status,
			title,
			visibility,
			join_code,
			state_json,
			version,
			fox_escape_at,
			created_by
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
		returning
			id,
			status,
			state_json,
			version,
			fox_escape_at,
			culprit_id,
			created_by,
			title,
			visibility,
			join_code
	`,
		"waiting",
		title,
		visibility,
		joinCodeValue,
		stateJSON,
		1,
		15,
		createdBy,
	).Scan(
		&row.ID,
		&row.Status,
		&row.StateJSON,
		&row.Version,
		&row.FoxEscapeAt,
		&row.CulpritID,
		&row.CreatedBy,
		&row.Title,
		&row.Visibility,
		&row.JoinCode,
	)
	if err != nil {
		return repo.GameRow{}, err
	}

	return row, nil
}

func (r *Repo) AddPlayer(ctx context.Context, tx *sql.Tx, gameID string, userID string, seat int) error {
	_, err := tx.ExecContext(ctx, `
		insert into game_players (game_id, user_id, seat)
		values ($1, $2, $3)
	`, gameID, userID, seat)
	return err
}

func (r *Repo) GetPlayers(ctx context.Context, gameID string) ([]repo.GamePlayerRow, error) {
	rows, err := r.db.QueryContext(ctx, `
		select gp.user_id, u.username, gp.seat
		from game_players gp
		join users u on u.id = gp.user_id
		where gp.game_id = $1
		order by gp.seat asc
	`, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPlayers(rows)
}

func (r *Repo) GetPlayersForUpdate(ctx context.Context, tx *sql.Tx, gameID string) ([]repo.GamePlayerRow, error) {
	rows, err := tx.QueryContext(ctx, `
		select gp.user_id, u.username, gp.seat
		from game_players gp
		join users u on u.id = gp.user_id
		where gp.game_id = $1
		order by gp.seat asc
	`, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPlayers(rows)
}

func (r *Repo) FindUnfinishedGameForUser(ctx context.Context, tx *sql.Tx, userID string) (repo.GameRow, error) {
	var row repo.GameRow
	err := tx.QueryRowContext(ctx, `
		select
			g.id,
			g.status,
			g.state_json,
			g.version,
			g.fox_escape_at,
			g.culprit_id,
			g.created_by,
			g.title,
			g.visibility,
			g.join_code
		from games g
		join game_players gp on gp.game_id = g.id
		where gp.user_id = $1
		  and g.status in ('waiting', 'active')
		order by g.created_at asc
		limit 1
		for update of g
	`, userID).Scan(
		&row.ID,
		&row.Status,
		&row.StateJSON,
		&row.Version,
		&row.FoxEscapeAt,
		&row.CulpritID,
		&row.CreatedBy,
		&row.Title,
		&row.Visibility,
		&row.JoinCode,
	)
	if err != nil {
		return repo.GameRow{}, err
	}
	return row, nil
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
		select
			id,
			status,
			state_json,
			version,
			fox_escape_at,
			culprit_id,
			created_by,
			title,
			visibility,
			join_code
		from games
		where id = $1
		for update
	`, gameID).Scan(
		&row.ID,
		&row.Status,
		&row.StateJSON,
		&row.Version,
		&row.FoxEscapeAt,
		&row.CulpritID,
		&row.CreatedBy,
		&row.Title,
		&row.Visibility,
		&row.JoinCode,
	)
	if err != nil {
		return repo.GameRow{}, err
	}
	return row, nil
}

func (r *Repo) GetGame(ctx context.Context, gameID string) (repo.GameRow, error) {
	var row repo.GameRow
	err := r.db.QueryRowContext(ctx, `
		select
			id,
			status,
			state_json,
			version,
			fox_escape_at,
			culprit_id,
			created_by,
			title,
			visibility,
			join_code
		from games
		where id = $1
	`, gameID).Scan(
		&row.ID,
		&row.Status,
		&row.StateJSON,
		&row.Version,
		&row.FoxEscapeAt,
		&row.CulpritID,
		&row.CreatedBy,
		&row.Title,
		&row.Visibility,
		&row.JoinCode,
	)
	if err != nil {
		return repo.GameRow{}, err
	}
	return row, nil
}

func (r *Repo) FindGameByJoinCodeForUpdate(ctx context.Context, tx *sql.Tx, code string) (repo.GameRow, error) {
	var row repo.GameRow
	err := tx.QueryRowContext(ctx, `
		select
			id,
			status,
			state_json,
			version,
			fox_escape_at,
			culprit_id,
			created_by,
			title,
			visibility,
			join_code
		from games
		where join_code = $1
		for update
	`, code).Scan(
		&row.ID,
		&row.Status,
		&row.StateJSON,
		&row.Version,
		&row.FoxEscapeAt,
		&row.CulpritID,
		&row.CreatedBy,
		&row.Title,
		&row.Visibility,
		&row.JoinCode,
	)
	if err != nil {
		return repo.GameRow{}, err
	}
	return row, nil
}

func (r *Repo) ListPublicWaitingGames(ctx context.Context) ([]repo.PublicGameRow, error) {
	rows, err := r.db.QueryContext(ctx, `
		select
			g.id,
			g.title,
			host.username as host_username,
			g.status,
			count(gp.user_id)::int as players_count
		from games g
		left join game_players gp on gp.game_id = g.id
		left join users host on host.id = g.created_by
		where g.status = 'waiting'
		  and g.visibility = 'public'
		group by g.id, g.title, host.username, g.status
		order by g.created_at desc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]repo.PublicGameRow, 0)
	for rows.Next() {
		var row repo.PublicGameRow
		if err := rows.Scan(
			&row.ID,
			&row.Title,
			&row.HostUsername,
			&row.Status,
			&row.PlayersCount,
		); err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repo) UpdateState(ctx context.Context, tx *sql.Tx, gameID string, status string, newStateJSON []byte, newVersion int) error {
	res, err := tx.ExecContext(ctx, `
		update games
		set status = $2,
		    state_json = $3,
		    version = $4,
		    updated_at = now()
		where id = $1
	`, gameID, status, newStateJSON, newVersion)
	if err != nil {
		return err
	}

	aff, _ := res.RowsAffected()
	if aff != 1 {
		return fmt.Errorf("update affected %d rows", aff)
	}
	return nil
}

func (r *Repo) RemovePlayer(ctx context.Context, tx *sql.Tx, gameID string, userID string) error {
	res, err := tx.ExecContext(ctx, `
		delete from game_players
		where game_id = $1 and user_id = $2
	`, gameID, userID)
	if err != nil {
		return err
	}

	aff, _ := res.RowsAffected()
	if aff != 1 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repo) SetGameCreator(ctx context.Context, tx *sql.Tx, gameID string, userID string) error {
	res, err := tx.ExecContext(ctx, `
		update games
		set created_by = $2,
		    updated_at = now()
		where id = $1
	`, gameID, userID)
	if err != nil {
		return err
	}

	aff, _ := res.RowsAffected()
	if aff != 1 {
		return fmt.Errorf("update affected %d rows", aff)
	}
	return nil
}

func (r *Repo) DeleteGame(ctx context.Context, tx *sql.Tx, gameID string) error {
	res, err := tx.ExecContext(ctx, `
		delete from games
		where id = $1
	`, gameID)
	if err != nil {
		return err
	}

	aff, _ := res.RowsAffected()
	if aff != 1 {
		return fmt.Errorf("delete affected %d rows", aff)
	}
	return nil
}

func scanPlayers(rows *sql.Rows) ([]repo.GamePlayerRow, error) {
	players := make([]repo.GamePlayerRow, 0)
	for rows.Next() {
		var player repo.GamePlayerRow
		if err := rows.Scan(&player.UserID, &player.Username, &player.Seat); err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return players, nil
}
