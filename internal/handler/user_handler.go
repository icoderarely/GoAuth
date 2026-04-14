package handler

import (
	"net/http"

	"github.com/icoderarely/GoAuth/internal/middleware"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	writeJSON(w, 200, map[string]any{
		"user_id":  claims.UserID,
		"username": claims.Username,
		"role":     claims.Role,
	})
}

func (h *UserHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	writeJSON(w, 200, map[string]string{
		"message": "welcome " + claims.Username,
	})
}

func (h *UserHandler) Admin(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	writeJSON(w, 200, map[string]string{
		"message": "admin access granted",
		"user":    claims.Username,
	})
}
