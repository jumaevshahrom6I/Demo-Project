package post

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"demo-project/internal/models"
	"demo-project/internal/repository"
	"demo-project/internal/service"
)

var _ service.PostService = (*Service)(nil)

const (
	maxTitleLen   = 200
	maxContentLen = 10_000
)

type Service struct {
	posts repository.PostRepository
}

func NewService(posts repository.PostRepository) *Service {
	return &Service{posts: posts}
}

// getOwnedPost находит пост и проверяет, что он принадлежит userID.
// Нет поста: models.ErrNotFound. Чужой пост: models.ErrForbidden.
// Общая для Update и Delete, чтобы не дублировать проверку.
func (s *Service) getOwnedPost(ctx context.Context, userID, id int64) (*models.Post, error) {
	p, err := s.posts.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get post: %w", err)
	}
	if p.UserID != userID {
		return nil, models.ErrForbidden
	}
	return p, nil
}

// validatePostInput обрезает пробелы и проверяет длины. Общая для создания и редактирования.
func validatePostInput(in models.PostInput) (models.PostInput, error) {
	in.Title = strings.TrimSpace(in.Title)
	in.Content = strings.TrimSpace(in.Content)

	if n := utf8.RuneCountInString(in.Title); n == 0 || n > maxTitleLen {
		return in, fmt.Errorf("%w: title must be 1 to %d characters", models.ErrInvalidInput, maxTitleLen)
	}
	if n := utf8.RuneCountInString(in.Content); n == 0 || n > maxContentLen {
		return in, fmt.Errorf("%w: content must be 1 to %d characters", models.ErrInvalidInput, maxContentLen)
	}
	return in, nil
}
