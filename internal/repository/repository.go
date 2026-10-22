package repository

import (
	"context"
	"time"

	"demo-project/internal/models"
)

type Repositories struct {
	User    UserRepository
	Role    RoleRepository
	Session SessionRepository
	Post    PostRepository
}

type UserRepository interface {
	Create(ctx context.Context, u *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type RoleRepository interface {
	GetByCode(ctx context.Context, code string) (*models.Role, error)
}

// Сессия — одна пара токенов: её id лежит в jti и access, и refresh-токена.
// Пока сессия есть в хранилище, оба токена действуют; удалили — оба отозваны.
type SessionRepository interface {
	Create(ctx context.Context, sessionID string, userID int64, ttl time.Duration) error
	Exists(ctx context.Context, sessionID string) (bool, error)
	Delete(ctx context.Context, sessionID string) error
}

type PostRepository interface {
	Create(ctx context.Context, p *models.Post) error
	List(ctx context.Context) ([]*models.Post, error)
	GetByID(ctx context.Context, id int64) (*models.Post, error)
	Update(ctx context.Context, p *models.Post) error
	Delete(ctx context.Context, id int64) error
}
