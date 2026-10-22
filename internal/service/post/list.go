package post

import (
	"context"
	"fmt"

	"demo-project/internal/models"
)

func (s *Service) List(ctx context.Context) ([]*models.Post, error) {
	posts, err := s.posts.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list posts: %w", err)
	}
	return posts, nil
}
