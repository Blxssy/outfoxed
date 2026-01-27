package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Level  string // debug/info/warn/error
	Pretty bool   // true = красивый консольный вывод
}

func Init(cfg Config) zerolog.Logger {
	var w io.Writer = os.Stdout

	if cfg.Pretty {
		w = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}

	log.Logger = zerolog.New(w).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()

	return log.Logger
}
