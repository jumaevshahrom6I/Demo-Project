package session

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"demo-project/internal/models"
)

type Repository struct {
	client *redis.Client
}

func NewRepository(client *redis.Client) *Repository {
	return &Repository{client: client}
}

func key(sessionID string) string {
	return "session:" + sessionID
}

func (r *Repository) Create(ctx context.Context, sessionID string, userID int64, ttl time.Duration) error {
	if err := r.client.Set(ctx, key(sessionID), strconv.FormatInt(userID, 10), ttl).Err(); err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

func (r *Repository) Exists(ctx context.Context, sessionID string) (bool, error) {
	n, err := r.client.Exists(ctx, key(sessionID)).Result()
	if err != nil {
		return false, fmt.Errorf("check session: %w", err)
	}
	return n > 0, nil
}

func (r *Repository) Delete(ctx context.Context, sessionID string) error {
	// DEL атомарен: из двух одновременных запросов с одной сессией удалит её только один.
	deleted, err := r.client.Del(ctx, key(sessionID)).Result()
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	if deleted == 0 {
		return models.ErrNotFound
	}
	return nil
}
