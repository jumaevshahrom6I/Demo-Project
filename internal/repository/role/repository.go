package role

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

func (r *Repository) GetByCode(ctx context.Context, code string) (*models.Role, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, code, description, created_at
		 FROM roles
		 WHERE code = $1`,
		code,
	)
	if err != nil {
		return nil, fmt.Errorf("select role by code: %w", err)
	}
	m, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("select role by code: %w", err)
	}
	return m.toModel(), nil
}
