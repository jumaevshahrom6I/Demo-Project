package post

import (
	"context"
	"fmt"
)

func (s *Service) Delete(ctx context.Context, userID, id int64) error {
	if _, err := s.getOwnedPost(ctx, userID, id); err != nil {
		return err
	}
	if err := s.posts.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
	return nil
}
