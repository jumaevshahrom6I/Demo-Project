package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"demo-project/internal/handler/response"
	"demo-project/internal/models"
	"demo-project/internal/service"
)

type claimsKey struct{}

func Auth(auth service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			if !ok || token == "" {
				response.Error(w, r, http.StatusUnauthorized, "missing bearer token", errors.New("no bearer token in Authorization header"))
				return
			}

			claims, err := auth.Authenticate(r.Context(), token)
			if err != nil {
				response.ServiceError(w, r, err)
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Вызывать только в обработчиках за middleware Auth: там claims точно есть.
func ClaimsFromContext(ctx context.Context) *models.AccessClaims {
	claims, _ := ctx.Value(claimsKey{}).(*models.AccessClaims)
	return claims
}
