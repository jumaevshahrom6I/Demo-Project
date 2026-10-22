package v1

import "demo-project/internal/service"

// Handler — ручки API версии 1. Методы разложены по файлам сущностей: auth.go, user.go.
// Каждая ручка делает три шага: разобрать запрос → вызвать сервис → ответить.
type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{services: services}
}
