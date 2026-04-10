package app

import (
	"context"
	"fox/config"
	"fox/internal/modules/game/transport/ws"
	userPg "fox/internal/modules/auth/repo/postgres"
	"fox/internal/modules/auth/service"
	userhttp "fox/internal/modules/auth/transport/http"
	httptransport "fox/internal/transport/http"
	"fox/pkg/logger"
	pg "fox/pkg/postgres"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	// lobbyhttp "fox/internal/modules/lobby/transport/http"
	// gamews "fox/internal/modules/game/transport/ws"
	// gamesvc "fox/internal/modules/game/service"
	// gamerepo "fox/internal/modules/game/repo/postgres"
)

func Run(cfg *config.Config) {
	log := logger.Init(logger.Config{
		Level:  "debug",
		Pretty: true,
	})

	log.Info().Msg("starting app")

	postgresDB, err := pg.New(cfg.PostgresConfig)
	if err != nil {
		log.Error().Err(err).Str("op", "postgres.New").Msg("error connecting to PostgreSQL")
		return
	}
	log.Info().Msg("successfully connected to PostgreSQL")

	v, err := pg.RunMigrations(postgresDB.DB, cfg.PostgresConfig)
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

	userRepo := userPg.NewUserRepository(postgresDB)
	refreshTokenRepo := userPg.NewRefreshTokenRepository(postgresDB)
	tokenManager := service.NewTokenManager(cfg.JWTSecret)
	authService := service.NewService(userRepo, refreshTokenRepo, tokenManager)
	authHandler := userhttp.NewHandler(authService, tokenManager)

	// Заглушки
	// var authHandler http.Handler  // = authhttp.NewRouter(...)
	// var lobbyHandler http.Handler // = lobbyhttp.NewRouter(...)
	// var wsHandler http.Handler    // = wsHandler

	hub := ws.NewHub()
	_ = hub

	router := httptransport.NewRouter(httptransport.Deps{
		Auth: authHandler,
		// Lobby:  lobbyHandler,
		GameWS: wsHandler,

		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	srv := &http.Server{
		Addr:         cfg.HTTP.Addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info().Str("addr", srv.Addr).Msg("http server listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("http server crashed")
		}
	}()

	<-ctx.Done()
	log.Info().Msg("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("http server shutdown error")
	} else {
		log.Info().Msg("http server stopped")
	}

	if err = postgresDB.Close(); err != nil {
		log.Error().Err(err).Msg("db close error")
	}
}
