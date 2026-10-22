package models

import "time"

type Post struct {
	ID        int64
	UserID    int64
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostInput struct {
	Title   string
	Content string
}
