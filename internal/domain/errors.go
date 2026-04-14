package domain

import (
	"errors"
)

// sentinel errors — return these from service, map to HTTP codes in handler
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("username already taken")
	ErrInvalidLogin      = errors.New("invalid login")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInvalidToken      = errors.New("invalid token")
)
