package handler

import (
	"encoding/json"
	"net/http"

	"github.com/icoderarely/GoAuth/internal/service"
)

type AuthHandler struct {
	auth service.AuthService
}

func NewAuthHandler(auth service.AuthService) *AuthHandler {
	return &AuthHandler{
		auth: auth,
	}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterRequest
	if err := readJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	user, err := h.auth.Register(r.Context(), payload.Username, payload.Password)
	if err != nil {
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user": user,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload LoginRequest

	if err := readJSON(r, &payload); err != nil {
		writeError(w, 400, "invalid request")
		return
	}

	token, err := h.auth.Login(r.Context(), payload.Username, payload.Password)
	if err != nil {
		writeError(w, 401, "invalid credentials")
		return
	}

	writeJSON(w, 200, map[string]string{
		"token": token,
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := map[string]any{
		"data": data,
	}

	return json.NewEncoder(w).Encode(resp)
}

func readJSON(r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(nil, r.Body, 1<<20) // 1MB

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	return dec.Decode(dst)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{
		"error": msg,
	})
}
