package v1

import (
	"net/http"

	"demo-project/internal/handler/middleware"
	"demo-project/internal/handler/response"
)

// GetMe работает только за middleware.Auth: id пользователя берётся из его access-токена.
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())

	user, err := h.services.User.GetByID(r.Context(), claims.UserID)
	if err != nil {
		response.ServiceError(w, r, err)
		return
	}

	response.JSON(w, http.StatusOK, newUserResponse(user))
}
