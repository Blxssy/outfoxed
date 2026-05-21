package repo

import (
	"context"
	"database/sql"
	"time"
)

type GameRow struct {
	ID          string
	Status      string
	StateJSON   []byte
	Version     int
	FoxEscapeAt int
	CulpritID   int
	CreatedBy   sql.NullString

	Title          string
	Visibility     string
	JoinCode       sql.NullString
	TurnDeadlineAt sql.NullTime
}

type GamePlayerRow struct {
	UserID   string
	Username string
	Seat     int
}

type PublicGameRow struct {
	ID           string
	Title        string
	HostUsername sql.NullString
	Status       string
	PlayersCount int
}

type Tx interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type GameRepo interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)

	CreateGame(
		ctx context.Context,
		tx *sql.Tx,
		createdBy string,
		title string,
		visibility string,
		joinCode *string,
		stateJSON []byte,
	) (GameRow, error)

	AddPlayer(ctx context.Context, tx *sql.Tx, gameID string, userID string, seat int) error
	GetPlayers(ctx context.Context, gameID string) ([]GamePlayerRow, error)
	GetPlayersForUpdate(ctx context.Context, tx *sql.Tx, gameID string) ([]GamePlayerRow, error)
	FindUnfinishedGameForUser(ctx context.Context, tx *sql.Tx, userID string) (GameRow, error)

	// Загружаем игру и блокируем строку до конца транзакции
	GetGameForUpdate(ctx context.Context, tx *sql.Tx, gameID string) (GameRow, error)

	// Загружаем игру без блокировки для read-only сценариев.
	GetGame(ctx context.Context, gameID string) (GameRow, error)

	// Проверяем, что пользователь участник игры
	IsPlayerInGame(ctx context.Context, tx *sql.Tx, gameID string, userID string) (bool, error)

	// Проверяем доступ к игре без транзакции.
	IsPlayerInGameReadonly(ctx context.Context, gameID string, userID string) (bool, error)

	// Сохраняем новый state и версию
	UpdateState(ctx context.Context, tx *sql.Tx, gameID string, status string, newStateJSON []byte, newVersion int) error

	ListPublicWaitingGames(ctx context.Context) ([]PublicGameRow, error)
	FindGameByJoinCodeForUpdate(ctx context.Context, tx *sql.Tx, code string) (GameRow, error)
	RemovePlayer(ctx context.Context, tx *sql.Tx, gameID string, userID string) error
	SetGameCreator(ctx context.Context, tx *sql.Tx, gameID string, userID string) error
	DeleteGame(ctx context.Context, tx *sql.Tx, gameID string) error
	UpdateStateAndDeadline(ctx context.Context, tx *sql.Tx, gameID string, status string, newStateJSON []byte, newVersion int, deadline *time.Time) error
	ListDueGamesForTimeout(ctx context.Context, limit int) ([]GameRow, error)
}
