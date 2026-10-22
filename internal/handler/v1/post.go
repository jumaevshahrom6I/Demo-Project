package v1

import (
	"net/http"

	"demo-project/internal/handler/middleware"
	"demo-project/internal/handler/response"
)

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req postRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	claims := middleware.ClaimsFromContext(r.Context())

	post, err := h.services.Post.Create(r.Context(), claims.UserID, req.toInput())
	if err != nil {
		response.ServiceError(w, r, err)
		return
	}

	response.JSON(w, http.StatusCreated, newPostResponse(post))
}

func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.services.Post.List(r.Context())
	if err != nil {
		response.ServiceError(w, r, err)
		return
	}

	// Пустой срез, а не nil: клиент получит [], а не null.
	resp := make([]postResponse, 0, len(posts))
	for _, p := range posts {
		resp = append(resp, newPostResponse(p))
	}

	response.JSON(w, http.StatusOK, resp)
}

func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	post, err := h.services.Post.GetByID(r.Context(), id)
	if err != nil {
		response.ServiceError(w, r, err)
		return
	}

	response.JSON(w, http.StatusOK, newPostResponse(post))
}

func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var req postRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	claims := middleware.ClaimsFromContext(r.Context())

	post, err := h.services.Post.Update(r.Context(), claims.UserID, id, req.toInput())
	if err != nil {
		response.ServiceError(w, r, err)
		return
	}

	response.JSON(w, http.StatusOK, newPostResponse(post))
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	claims := middleware.ClaimsFromContext(r.Context())

	if err := h.services.Post.Delete(r.Context(), claims.UserID, id); err != nil {
		response.ServiceError(w, r, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
