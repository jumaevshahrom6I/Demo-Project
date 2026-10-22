package v1

import (
	"time"

	"demo-project/internal/models"
)

type postRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (req postRequest) toInput() models.PostInput {
	return models.PostInput{
		Title:   req.Title,
		Content: req.Content,
	}
}

type postResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newPostResponse(p *models.Post) postResponse {
	return postResponse{
		ID:        p.ID,
		UserID:    p.UserID,
		Title:     p.Title,
		Content:   p.Content,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
