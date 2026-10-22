package auth

import (
	"context"
	"errors"
	"fmt"

	"demo-project/internal/models"
)

// Logout удаляет сессию: и access, и refresh-токен этой пары сразу перестают действовать.
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	claims, err := s.parse(refreshToken, refreshTokenType)
	if err != nil {
		return err
	}

	// Сессии уже нет — пользователь и так вышел, это не ошибка.
	if err := s.sessions.Delete(ctx, claims.ID); err != nil && !errors.Is(err, models.ErrNotFound) {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}
