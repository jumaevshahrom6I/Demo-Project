package post

import (
	"context"
	"fmt"

	"demo-project/internal/models"
)

func (s *Service) Create(ctx context.Context, userID int64, in models.PostInput) (*models.Post, error) {
	in, err := validatePostInput(in)
	if err != nil {
		return nil, err
	}

	p := &models.Post{
		UserID:  userID,
		Title:   in.Title,
		Content: in.Content,
	}
	if err := s.posts.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}
	return p, nil
}
