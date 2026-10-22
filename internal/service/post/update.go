package post

import (
	"context"
	"fmt"

	"demo-project/internal/models"
)

func (s *Service) Update(ctx context.Context, userID, id int64, in models.PostInput) (*models.Post, error) {
	// Порядок проверок: пост есть (404) → автор (403) → валидация (400) → запись.
	p, err := s.getOwnedPost(ctx, userID, id)
	if err != nil {
		return nil, err
	}

	in, err = validatePostInput(in)
	if err != nil {
		return nil, err
	}

	p.Title = in.Title
	p.Content = in.Content
	if err := s.posts.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("update post: %w", err)
	}
	return p, nil
}
