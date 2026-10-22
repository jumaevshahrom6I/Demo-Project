package v1

import (
	"time"

	"demo-project/internal/models"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

func (req registerRequest) toInput() models.RegisterInput {
	return models.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		RoleCode: req.Role,
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req loginRequest) toInput() models.LoginInput {
	return models.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type tokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func newTokenPairResponse(pair *models.TokenPair) tokenPairResponse {
	return tokenPairResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	}
}

// userResponse — пользователь в ответе API
type userResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	RoleID    int64     `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(u *models.User) userResponse {
	return userResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		RoleID:    u.RoleID,
		CreatedAt: u.CreatedAt,
	}
}
