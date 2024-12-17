package middleware

import (
	"api-gateway/internal/api/proto/auth"
	"api-gateway/internal/clients"
	"context"
	"net/http"
	"strings"
)

// AuthMiddleware возвращает функцию-обёртку для обработчиков, добавляя проверку токена
func AuthMiddleware(authClient *clients.AuthClient) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Извлечение токена из заголовка Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			// Удаляем "Bearer " из заголовка
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == authHeader {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			// Проверяем токен через AuthService
			validateReq := &auth.ValidateTokenRequest{Token: token}
			res, err := authClient.ValidateToken(context.Background(), validateReq)
			if err != nil || !res.Valid {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Добавляем информацию о роли в контекст
			ctx := context.WithValue(r.Context(), "role", res.Role)

			// Передаём управление следующему обработчику
			next(w, r.WithContext(ctx))
		}
	}
}
