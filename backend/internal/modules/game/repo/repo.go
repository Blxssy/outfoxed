package repo

import (
	"context"
	"database/sql"
)

type GameRow struct {
	ID          string
	Status      string
	StateJSON   []byte
	Version     int
	FoxEscapeAt int
	CulpritID   int
}

type Tx interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type GameRepo interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)

	// Загружаем игру и блокируем строку до конца транзакции
	GetGameForUpdate(ctx context.Context, tx *sql.Tx, gameID string) (GameRow, error)

	// Проверяем, что пользователь участник игры
	IsPlayerInGame(ctx context.Context, tx *sql.Tx, gameID string, userID string) (bool, error)

	// Сохраняем новый state и версию
	UpdateState(ctx context.Context, tx *sql.Tx, gameID string, newStateJSON []byte, newVersion int) error
}
