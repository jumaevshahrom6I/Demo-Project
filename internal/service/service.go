package service

import (
	"context"

	"demo-project/internal/models"
)

// Services — все сервисы приложения. Собирается в app, передаётся в handler.
type Services struct {
	Auth AuthService
	User UserService
	Post PostService
}

type AuthService interface {
	Register(ctx context.Context, in models.RegisterInput) (*models.User, error)
	Login(ctx context.Context, in models.LoginInput) (*models.TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (*models.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	Authenticate(ctx context.Context, accessToken string) (*models.AccessClaims, error)
}

type UserService interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

type PostService interface {
	Create(ctx context.Context, userID int64, in models.PostInput) (*models.Post, error)
	List(ctx context.Context) ([]*models.Post, error)
	GetByID(ctx context.Context, id int64) (*models.Post, error)
	Update(ctx context.Context, userID, id int64, in models.PostInput) (*models.Post, error)
	Delete(ctx context.Context, userID, id int64) error
}
