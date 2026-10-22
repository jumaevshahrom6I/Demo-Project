package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"demo-project/internal/config"
	"demo-project/internal/handler"
	"demo-project/internal/repository"
	postrepo "demo-project/internal/repository/post"
	rolerepo "demo-project/internal/repository/role"
	sessionrepo "demo-project/internal/repository/session"
	userrepo "demo-project/internal/repository/user"
	"demo-project/internal/service"
	authservice "demo-project/internal/service/auth"
	postservice "demo-project/internal/service/post"
	userservice "demo-project/internal/service/user"
	"demo-project/pkg/logger"
	"demo-project/pkg/postgres"
	"demo-project/pkg/redis"
)

func Run() (err error) {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log, closeLog, err := logger.New(logger.Config{
		Env:   cfg.App.Env,
		Level: cfg.App.LogLevel,
		File:  cfg.App.LogFile,
	})
	if err != nil {
		return err
	}
	defer func() { _ = closeLog() }()
	slog.SetDefault(log)

	// Пишем причину остановки, пока файл лога ещё открыт: defer выполняются в обратном порядке.
	defer func() {
		if err != nil {
			slog.Error("app stopped", "error", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	postgresPool, err := postgres.New(ctx, postgres.Config{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		Database: cfg.Postgres.Database,
		SSLMode:  cfg.Postgres.SSLMode,
	})
	if err != nil {
		return err
	}
	defer postgresPool.Close()

	redisClient, err := redis.New(ctx, redis.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		return err
	}
	defer func() { _ = redisClient.Close() }()

	repos := &repository.Repositories{
		User:    userrepo.NewRepository(postgresPool),
		Role:    rolerepo.NewRepository(postgresPool),
		Session: sessionrepo.NewRepository(redisClient),
		Post:    postrepo.NewRepository(postgresPool),
	}

	services := &service.Services{
		Auth: authservice.NewService(repos.User, repos.Role, repos.Session, authservice.TokenConfig{
			Secret:     cfg.JWT.Secret,
			AccessTTL:  cfg.JWT.AccessTTL,
			RefreshTTL: cfg.JWT.RefreshTTL,
		}),
		User: userservice.NewService(repos.User),
		Post: postservice.NewService(repos.Post),
	}

	handlers := handler.NewHandler(services)

	server := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      handlers.InitRoutes(),
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	return serve(ctx, server, cfg.HTTP.ShutdownTimeout)
}

func serve(ctx context.Context, server *http.Server, shutdownTimeout time.Duration) error {
	serverErr := make(chan error, 1)
	go func() {
		slog.Info("server started", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("listen and serve: %w", err)
	case <-ctx.Done():
	}

	slog.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}
	return nil
}
