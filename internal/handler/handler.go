package handler

import (
	"net/http"

	"demo-project/internal/handler/middleware"
	v1 "demo-project/internal/handler/v1"
	"demo-project/internal/service"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{services: services}
}

// InitRoutes регистрирует все маршруты API и оборачивает их в общие middleware.
func (h *Handler) InitRoutes() http.Handler {
	apiV1 := v1.NewHandler(h.services)
	requireAuth := middleware.Auth(h.services.Auth)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/auth/register", apiV1.Register)
	mux.HandleFunc("POST /api/v1/auth/login", apiV1.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", apiV1.Refresh)
	mux.HandleFunc("POST /api/v1/auth/logout", apiV1.Logout)

	mux.Handle("GET /api/v1/users/me", requireAuth(http.HandlerFunc(apiV1.GetMe)))

	mux.Handle("POST /api/v1/posts", requireAuth(http.HandlerFunc(apiV1.CreatePost)))
	mux.HandleFunc("GET /api/v1/posts", apiV1.ListPosts)
	mux.HandleFunc("GET /api/v1/posts/{id}", apiV1.GetPost)
	mux.Handle("PUT /api/v1/posts/{id}", requireAuth(http.HandlerFunc(apiV1.UpdatePost)))
	mux.Handle("DELETE /api/v1/posts/{id}", requireAuth(http.HandlerFunc(apiV1.DeletePost)))

	// Logging снаружи: он видит и ответ 500, который отдаст Recover после паники.
	return middleware.CORS(middleware.Logging(middleware.Recover(mux)))
}
