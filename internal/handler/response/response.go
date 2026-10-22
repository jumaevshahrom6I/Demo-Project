package response

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"demo-project/internal/models"
)

type errorBody struct {
	Error string `json:"error"`
}

func JSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if body != nil {
		_ = json.NewEncoder(w).Encode(body)
	}
}

// cause видит только лог, клиенту уходит message.
func Error(w http.ResponseWriter, r *http.Request, status int, message string, cause error) {
	level := slog.LevelWarn
	if status >= http.StatusInternalServerError {
		level = slog.LevelError
	}
	slog.Log(r.Context(), level, "request failed",
		"method", r.Method,
		"path", r.URL.Path,
		"status", status,
		"error", cause,
	)
	JSON(w, status, errorBody{Error: message})
}

func ServiceError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, models.ErrInvalidInput):
		// Текст собран в сервисе для клиента: "invalid input: invalid email".
		Error(w, r, http.StatusBadRequest, err.Error(), err)
	case errors.Is(err, models.ErrInvalidCredentials):
		Error(w, r, http.StatusUnauthorized, models.ErrInvalidCredentials.Error(), err)
	case errors.Is(err, models.ErrInvalidToken):
		Error(w, r, http.StatusUnauthorized, models.ErrInvalidToken.Error(), err)
	case errors.Is(err, models.ErrEmailTaken):
		Error(w, r, http.StatusConflict, models.ErrEmailTaken.Error(), err)
	case errors.Is(err, models.ErrForbidden):
		Error(w, r, http.StatusForbidden, models.ErrForbidden.Error(), err)
	case errors.Is(err, models.ErrNotFound):
		Error(w, r, http.StatusNotFound, models.ErrNotFound.Error(), err)
	default:
		Error(w, r, http.StatusInternalServerError, "internal server error", err)
	}
}
