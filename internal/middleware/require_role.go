package middleware

import (
	"net/http"

	"github.com/icoderarely/GoAuth/internal/domain"
	"github.com/icoderarely/GoAuth/internal/service"
)

func RequireRole(required domain.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims, ok := r.Context().Value(claimsKey).(*service.Claims)
			if !ok || claims == nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if claims.Role != required {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
