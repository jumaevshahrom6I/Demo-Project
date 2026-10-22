package config

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App      App
	HTTP     HTTP
	Postgres Postgres
	Redis    Redis
	JWT      JWT
}

type App struct {
	Env      string     `envconfig:"APP_ENV" default:"local"` // local | prod
	LogLevel slog.Level `envconfig:"LOG_LEVEL" default:"info"`
	LogFile  string     `envconfig:"LOG_FILE" default:"app.log"` // только для APP_ENV=local
}

type HTTP struct {
	Port            string        `envconfig:"HTTP_PORT" default:"8080"`
	ReadTimeout     time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"10s"`
	WriteTimeout    time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"10s"`
	ShutdownTimeout time.Duration `envconfig:"HTTP_SHUTDOWN_TIMEOUT" default:"10s"`
}

type Postgres struct {
	Host     string `envconfig:"POSTGRES_HOST" default:"localhost"`
	Port     string `envconfig:"POSTGRES_PORT" default:"5432"`
	User     string `envconfig:"POSTGRES_USER" default:"demo"`
	Password string `envconfig:"POSTGRES_PASSWORD" default:"demo"`
	Database string `envconfig:"POSTGRES_DB" default:"demo-db"`
	SSLMode  string `envconfig:"POSTGRES_SSLMODE" default:"disable"`
}

type Redis struct {
	Host     string `envconfig:"REDIS_HOST" default:"localhost"`
	Port     string `envconfig:"REDIS_PORT" default:"6379"`
	Password string `envconfig:"REDIS_PASSWORD" default:""`
	DB       int    `envconfig:"REDIS_DB" default:"0"`
}

type JWT struct {
	// Секрет без значения по умолчанию: подписывать токены известной всем строкой нельзя.
	Secret     string        `envconfig:"JWT_SECRET" required:"true"`
	AccessTTL  time.Duration `envconfig:"JWT_ACCESS_TTL" default:"15m"`
	RefreshTTL time.Duration `envconfig:"JWT_REFRESH_TTL" default:"720h"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("load .env: %w", err)
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return &cfg, nil
}
