package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/icoderarely/GoAuth/internal/service"
)

type ctxKey string

const claimsKey ctxKey = "claims"

func AuthMiddleware(auth service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "unauthorized", 401)
				return
			}

			token := parts[1]

			claims, err := auth.ValidateToken(token)
			if err != nil {
				http.Error(w, "unauthorized", 401)
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
