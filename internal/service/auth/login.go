package auth

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"demo-project/internal/models"
)

func (s *Service) Login(ctx context.Context, in models.LoginInput) (*models.TokenPair, error) {
	u, err := s.users.GetByEmail(ctx, normalizeEmail(in.Email))
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, models.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)); err != nil {
		return nil, models.ErrInvalidCredentials
	}

	pair, err := s.newTokenPair(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("create tokens: %w", err)
	}
	return pair, nil
}
