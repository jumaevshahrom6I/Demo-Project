package auth

import (
	"context"
	"fmt"

	"demo-project/internal/models"
)

// Authenticate проверяет access-токен: подпись, срок и то, что его сессию не отозвали (logout, refresh).
func (s *Service) Authenticate(ctx context.Context, accessToken string) (*models.AccessClaims, error) {
	claims, err := s.parse(accessToken, accessTokenType)
	if err != nil {
		return nil, err
	}
	userID, err := claims.userID()
	if err != nil {
		return nil, err
	}

	active, err := s.sessions.Exists(ctx, claims.ID)
	if err != nil {
		return nil, fmt.Errorf("check session: %w", err)
	}
	if !active {
		return nil, fmt.Errorf("%w: session revoked", models.ErrInvalidToken)
	}

	return &models.AccessClaims{UserID: userID, RoleID: claims.RoleID}, nil
}
