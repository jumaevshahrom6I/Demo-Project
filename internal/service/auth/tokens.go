package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"demo-project/internal/models"
)

const (
	accessTokenType  = "access"
	refreshTokenType = "refresh"
)

type tokenClaims struct {
	jwt.RegisteredClaims
	Type   string `json:"type"`
	RoleID int64  `json:"role_id,omitempty"`
}

// newTokenPair создаёт сессию и выпускает на неё пару токенов.
func (s *Service) newTokenPair(ctx context.Context, u *models.User) (*models.TokenPair, error) {
	sessionID := rand.Text()
	if err := s.sessions.Create(ctx, sessionID, u.ID, s.refreshTTL); err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	now := time.Now()

	accessToken, err := s.sign(tokenClaims{
		RegisteredClaims: registeredClaims(u.ID, sessionID, now, s.accessTTL),
		Type:             accessTokenType,
		RoleID:           u.RoleID,
	})
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	// Роль в refresh не кладём: при обновлении она берётся заново из БД.
	refreshToken, err := s.sign(tokenClaims{
		RegisteredClaims: registeredClaims(u.ID, sessionID, now, s.refreshTTL),
		Type:             refreshTokenType,
	})
	if err != nil {
		return nil, fmt.Errorf("sign refresh token: %w", err)
	}

	return &models.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func registeredClaims(userID int64, sessionID string, now time.Time, ttl time.Duration) jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10), // sub — стандартное поле JWT для id пользователя
		ID:        sessionID,                     // jti — у обоих токенов пары он один: id сессии
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}
}

func (s *Service) sign(claims tokenClaims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}

// parse проверяет подпись, срок и тип токена. Действует ли сессия, проверяет вызывающий.
func (s *Service) parse(token, wantType string) (*tokenClaims, error) {
	var claims tokenClaims
	_, err := jwt.ParseWithClaims(token, &claims,
		func(*jwt.Token) (any, error) { return s.jwtSecret, nil },
		// Принимаем только наш алгоритм: иначе можно подсунуть токен, подписанный по-другому.
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", models.ErrInvalidToken, err)
	}
	if claims.Type != wantType {
		return nil, fmt.Errorf("%w: expected %s token, got %q", models.ErrInvalidToken, wantType, claims.Type)
	}
	return &claims, nil
}

func (c *tokenClaims) userID() (int64, error) {
	id, err := strconv.ParseInt(c.Subject, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: invalid subject", models.ErrInvalidToken)
	}
	return id, nil
}
