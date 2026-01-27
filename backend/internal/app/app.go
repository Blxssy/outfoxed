package app

import (
	"fox/config"
	"fox/pkg/logger"
	"fox/pkg/postgres"
)

func Run(cfg *config.Config) {
	log := logger.Init(logger.Config{
		Level:  "debug",
		Pretty: true,
	})

	log.Info().Msg("starting app")

	postgresDB, err := postgres.New(cfg.PostgresConfig)
	if err != nil {
		log.Error().Err(err).Str("op", "postgres.New").Msg("error connecting to PostgreSQL")
		return
	}
	log.Info().Msg("successfully connected to PostgreSQL")

	v, err := postgres.RunMigrations(postgresDB.DB, cfg.PostgresConfig)
	if err != nil {
		log.Error().
			Err(err).
			Str("op", "postgres.RunMigrations").
			Msg("error running PostgreSQL migrations")
		return
	}
	log.Info().
		Uint("version", v).
		Msg("successfully completed PostgreSQL migrations")
}
