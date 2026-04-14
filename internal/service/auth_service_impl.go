package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/icoderarely/GoAuth/internal/domain"
	"github.com/icoderarely/GoAuth/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID               string
	Username             string
	Role                 domain.Role
	jwt.RegisteredClaims // embeds ExpiresAt, IssuedAt, etc.
}

type AuthServiceImpl struct {
	userRepo  repository.UserRepository
	jwtSecret string
	tokenTTL  time.Duration
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string, tokenTTL time.Duration) *AuthServiceImpl {
	if userRepo == nil || jwtSecret == "" || tokenTTL <= 0 {
		return nil
	}

	return &AuthServiceImpl{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		tokenTTL:  tokenTTL,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, username, password string) (*domain.User, error) {
	// bcrypt hash + repo.Save

	if username == "" {
		return nil, domain.ErrInvalidInput
	}
	if password == "" {
		return nil, domain.ErrInvalidInput
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := NewUser(username, hash)

	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	out := *user
	out.PasswordHash = ""

	return &out, nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, username, password string) (token string, err error) {
	// FindByUsername + bcrypt.Compare + jwt.NewWithClaims + SignedString
	if username == "" {
		return "", domain.ErrInvalidInput
	}
	if password == "" {
		return "", domain.ErrInvalidInput
	}

	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", domain.ErrInvalidLogin
		}
		return "", err
	}
	if user == nil {
		return "", domain.ErrInvalidLogin
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", domain.ErrInvalidLogin
	}

	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(s.tokenTTL)),
			Subject:   user.ID,
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := jwtToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (s *AuthServiceImpl) ValidateToken(tokenStr string) (*Claims, error) {
	// jwt.ParseWithClaims
	if tokenStr == "" {
		return nil, domain.ErrInvalidToken
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(token *jwt.Token) (any, error) {
			method, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok || method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, domain.ErrInvalidToken
			}
			return []byte(s.jwtSecret), nil
		},
	)

	if err != nil || token == nil || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	return claims, nil
}

func NewUser(username string, hash []byte) *domain.User {
	return &domain.User{
		ID:           uuid.NewString(),
		Username:     username,
		PasswordHash: string(hash),
		Role:         domain.RoleUser,
		CreatedAt:    time.Now().UTC(),
	}
}
