package config

import (
	"log"
	"sync"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTP           HTTPConfig
	PostgresConfig PostgresConfig
	JWT
}

type HTTPConfig struct {
	Addr string `env:"HTTP_ADDR" envDefault:":8080"`
}

type PostgresConfig struct {
	DataSource string `env:"DB_DATA_SOURCE,required"`
}

type JWT struct {
	JWTSecret string `env:"JWT_SECRET,required"`
}

var (
	config Config
	once   sync.Once
)

func Get() *Config {
	once.Do(func() {
		_ = godotenv.Load()
		if err := env.Parse(&config); err != nil {
			log.Fatal(err)
		}
	})
	return &config
}
