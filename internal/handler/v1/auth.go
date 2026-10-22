package v1

import (
	"net/http"

	"demo-project/internal/handler/response"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	user, err := h.services.Auth.Register(r.Context(), req.toInput())
	if err != nil {
		response.ServiceError(w, r, err)
		return
	}

	response.JSON(w, http.StatusCreated, newUserResponse(user))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	tokens, err := h.services.Auth.Login(r.Context(), req.toInput())
	if err != nil {
		response.ServiceError(w, r, err)
		return
	}

	response.JSON(w, http.StatusOK, newTokenPairResponse(tokens))
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshTokenRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	tokens, err := h.services.Auth.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		response.ServiceError(w, r, err)
		return
	}

	response.JSON(w, http.StatusOK, newTokenPairResponse(tokens))
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var req refreshTokenRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	if err := h.services.Auth.Logout(r.Context(), req.RefreshToken); err != nil {
		response.ServiceError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
