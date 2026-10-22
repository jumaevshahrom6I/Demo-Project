package role

import (
	"time"

	"demo-project/internal/models"
)

type dbModel struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Code        string    `db:"code"`
	Description *string   `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

func (m dbModel) toModel() *models.Role {
	r := &models.Role{
		ID:   m.ID,
		Code: m.Code,
		Name: m.Name,
	}
	if m.Description != nil {
		r.Description = *m.Description
	}
	return r
}
