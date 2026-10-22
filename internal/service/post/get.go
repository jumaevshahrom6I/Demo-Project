package post

import (
	"context"
	"fmt"

	"demo-project/internal/models"
)

func (s *Service) GetByID(ctx context.Context, id int64) (*models.Post, error) {
	p, err := s.posts.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get post: %w", err)
	}
	return p, nil
}
