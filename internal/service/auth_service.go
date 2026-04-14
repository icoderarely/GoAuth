package service

import (
	"context"

	"github.com/icoderarely/GoAuth/internal/domain"
)

type AuthService interface {
	Register(ctx context.Context, username, password string) (*domain.User, error)
	Login(ctx context.Context, username, password string) (token string, err error)
	ValidateToken(tokenStr string) (*Claims, error)
}
