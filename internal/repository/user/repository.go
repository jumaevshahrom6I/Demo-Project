package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"demo-project/internal/models"
)

const uniqueViolation = "23505"

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, u *models.User) error {
	err := r.db.QueryRow(ctx,
		`INSERT INTO users (email, password_hash, name, role_id)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, created_at`,
		u.Email, u.PasswordHash, u.Name, u.RoleID,
	).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return models.ErrEmailTaken
		}
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, email, password_hash, name, role_id, created_at
		 FROM users
		 WHERE id = $1`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("select user by id: %w", err)
	}
	return collectUser(rows)
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, email, password_hash, name, role_id, created_at
		 FROM users
		 WHERE email = $1`,
		email,
	)
	if err != nil {
		return nil, fmt.Errorf("select user by email: %w", err)
	}
	return collectUser(rows)
}

func collectUser(rows pgx.Rows) (*models.User, error) {
	m, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("scan user: %w", err)
	}
	return m.toModel(), nil
}
