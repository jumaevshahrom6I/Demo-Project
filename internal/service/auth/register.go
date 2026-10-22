package auth

import (
	"context"
	"fmt"
	"net/mail"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"

	"demo-project/internal/models"
)

const (
	minPasswordLen   = 8
	maxPasswordBytes = 72 // bcrypt учитывает только первые 72 байта пароля
)

// Роли, которые можно выбрать при регистрации. Всё, чего здесь нет (admin, moderator, auditor...),
// отклоняется: служебные роли назначает только администратор.
var selfRegisterRoles = map[string]bool{
	"user": true,
}

func (s *Service) Register(ctx context.Context, in models.RegisterInput) (*models.User, error) {
	in.Email = normalizeEmail(in.Email)
	in.Name = strings.TrimSpace(in.Name)
	if err := validateRegister(in); err != nil {
		return nil, err
	}

	// Роль уже проверена по selfRegisterRoles: если её нет в БД, не применена миграция — это ошибка сервера.
	role, err := s.roles.GetByCode(ctx, in.RoleCode)
	if err != nil {
		return nil, fmt.Errorf("get role %q: %w", in.RoleCode, err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	u := &models.User{
		Email:        in.Email,
		PasswordHash: string(hash),
		Name:         in.Name,
		RoleID:       role.ID,
	}
	// Занятость email не проверяем заранее: это сделает UNIQUE в БД (models.ErrEmailTaken).
	if err := s.users.Create(ctx, u); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return u, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validateRegister(in models.RegisterInput) error {
	if addr, err := mail.ParseAddress(in.Email); err != nil || addr.Address != in.Email {
		return fmt.Errorf("%w: invalid email", models.ErrInvalidInput)
	}
	if in.Name == "" {
		return fmt.Errorf("%w: name is required", models.ErrInvalidInput)
	}
	if utf8.RuneCountInString(in.Password) < minPasswordLen {
		return fmt.Errorf("%w: password must be at least %d characters", models.ErrInvalidInput, minPasswordLen)
	}
	if len(in.Password) > maxPasswordBytes {
		return fmt.Errorf("%w: password must be at most %d bytes", models.ErrInvalidInput, maxPasswordBytes)
	}
	if !selfRegisterRoles[in.RoleCode] {
		return fmt.Errorf("%w: role %q is not available for registration", models.ErrInvalidInput, in.RoleCode)
	}
	return nil
}
