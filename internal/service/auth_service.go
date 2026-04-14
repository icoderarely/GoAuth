package service

import (
	"context"
)

type AuthService interface {
	Register(ctx context.Context, username, password string) (*UserResponse, error)
	Login(ctx context.Context, username, password string) (token string, err error)
	ValidateToken(tokenStr string) (*Claims, error)
}
