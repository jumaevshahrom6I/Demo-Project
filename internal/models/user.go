package models

import "time"

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	Name         string
	RoleID       int64
	CreatedAt    time.Time
}

type RegisterInput struct {
	Email    string
	Password string
	Name     string
	RoleCode string
}

type LoginInput struct {
	Email    string
	Password string
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type AccessClaims struct {
	UserID int64
	RoleID int64
}
