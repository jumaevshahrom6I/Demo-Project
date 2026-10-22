package auth

import (
	"context"
	"errors"
	"fmt"

	"demo-project/internal/models"
)

// Refresh закрывает старую сессию (её access и refresh больше не действуют) и выдаёт новую пару.
func (s *Service) Refresh(ctx context.Context, refreshToken string) (*models.TokenPair, error) {
	claims, err := s.parse(refreshToken, refreshTokenType)
	if err != nil {
		return nil, err
	}
	userID, err := claims.userID()
	if err != nil {
		return nil, err
	}

	// Сессии нет — этот refresh уже использовали или пользователь вышел.
	if err := s.sessions.Delete(ctx, claims.ID); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, models.ErrInvalidToken
		}
		return nil, fmt.Errorf("delete session: %w", err)
	}

	// Пользователя берём из БД, а не из токена: так в новый access попадёт актуальная роль.
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, models.ErrInvalidToken
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	pair, err := s.newTokenPair(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("create tokens: %w", err)
	}
	return pair, nil
}
