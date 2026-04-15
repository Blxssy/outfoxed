package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"fox/config"
	"fox/internal/repo/migrations"
	"time"

	"github.com/golang-migrate/migrate/v4/source/iofs"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	driverName    = "pgx"
	defaultSchema = "public"
)

func New(cfg config.PostgresConfig) (*sqlx.DB, error) {
	sqlxDB, err := sqlx.Open(driverName, cfg.DataSource)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Open: %w", err)
	}

	sqlxDB.SetMaxOpenConns(10)
	sqlxDB.SetMaxIdleConns(10 / 2)
	sqlxDB.SetConnMaxLifetime(5 * time.Minute)

	if err = sqlxDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	return sqlxDB, nil
}

func RunMigrations(instance *sql.DB, cfg config.PostgresConfig) (uint, error) {
	if _, err := instance.Exec("create schema if not exists " + defaultSchema); err != nil {
		return 0, fmt.Errorf("create schema: %w", err)
	}

	driver, err := postgres.WithInstance(instance, &postgres.Config{
		SchemaName: defaultSchema,
	})
	if err != nil {
		return 0, fmt.Errorf("create driver with instance: %w", err)
	}

	sub, err := fs.Sub(migrations.FS, ".")
	if err != nil {
		return 0, fmt.Errorf("fs.Sub migrations: %w", err)
	}

	src, err := iofs.New(sub, ".")
	if err != nil {
		return 0, fmt.Errorf("iofs.New: %w", err)
	}

	migrateInst, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		return 0, fmt.Errorf("create migrate instance: %w", err)
	}

	if err = migrateInst.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return 0, fmt.Errorf("up migrations: %w", err)
	}

	version, _, _ := migrateInst.Version()
	return version, nil
}
