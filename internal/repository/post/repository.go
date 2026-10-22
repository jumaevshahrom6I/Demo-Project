package post

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"demo-project/internal/models"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, p *models.Post) error {
	err := r.db.QueryRow(ctx,
		`INSERT INTO posts (user_id, title, content)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_at, updated_at`,
		p.UserID, p.Title, p.Content,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert post: %w", err)
	}
	return nil
}

// List возвращает все посты, новые первыми. Если постов нет, срез пустой, а не nil.
func (r *Repository) List(ctx context.Context) ([]*models.Post, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, title, content, created_at, updated_at
		 FROM posts
		 ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("select posts: %w", err)
	}
	defer rows.Close()

	posts := make([]*models.Post, 0)
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan post: %w", err)
		}
		posts = append(posts, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate posts: %w", err)
	}
	return posts, nil
}

// GetByID возвращает пост по id или models.ErrNotFound, если такого поста нет.
func (r *Repository) GetByID(ctx context.Context, id int64) (*models.Post, error) {
	var p models.Post
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, title, content, created_at, updated_at
		 FROM posts
		 WHERE id = $1`,
		id,
	).Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("select post: %w", err)
	}
	return &p, nil
}

// Update меняет title и content, обновляет updated_at и записывает его обратно в p.
func (r *Repository) Update(ctx context.Context, p *models.Post) error {
	err := r.db.QueryRow(ctx,
		`UPDATE posts
		 SET title = $1, content = $2, updated_at = now()
		 WHERE id = $3
		 RETURNING updated_at`,
		p.Title, p.Content, p.ID,
	).Scan(&p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ErrNotFound
		}
		return fmt.Errorf("update post: %w", err)
	}
	return nil
}

// Delete удаляет пост. Если строки не было, возвращает models.ErrNotFound.
func (r *Repository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM posts WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}
