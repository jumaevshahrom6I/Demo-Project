package auth

import (
	"time"

	"demo-project/internal/repository"
	"demo-project/internal/service"
)

// Компилятор проверит, что Service реализует интерфейс service.AuthService.
var _ service.AuthService = (*Service)(nil)

type TokenConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type Service struct {
	users      repository.UserRepository
	roles      repository.RoleRepository
	sessions   repository.SessionRepository
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewService(
	users repository.UserRepository,
	roles repository.RoleRepository,
	sessions repository.SessionRepository,
	tokens TokenConfig,
) *Service {
	return &Service{
		users:      users,
		roles:      roles,
		sessions:   sessions,
		jwtSecret:  []byte(tokens.Secret),
		accessTTL:  tokens.AccessTTL,
		refreshTTL: tokens.RefreshTTL,
	}
}
