package inmemory

import (
	"context"
	"sync"

	"github.com/icoderarely/GoAuth/internal/domain"
)

// In-memory implementation lives in repository/inmemory/user_store.go
// Uses a sync.RWMutex-protected map[string]*domain.User

type Store struct {
	users map[string]*domain.User
	mu    sync.RWMutex
}

func cloneUser(u *domain.User) *domain.User {
	if u == nil {
		return nil
	}
	c := *u // value copy of struct
	return &c
}

func NewStore(users map[string]*domain.User) *Store {
	if users == nil {
		users = make(map[string]*domain.User)
	}
	return &Store{users: users}
}

func (s *Store) Save(ctx context.Context, user *domain.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if user == nil {
		return domain.ErrInvalidInput
	}

	for _, u := range s.users {
		if u.Username == user.Username {
			return domain.ErrUserAlreadyExists
		}
	}

	s.users[user.ID] = cloneUser(user)

	return nil
}

func (s *Store) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if username == "" {
		return nil, domain.ErrInvalidInput
	}

	for _, u := range s.users {
		if u.Username == username {
			return cloneUser(u), nil
		}
	}

	return nil, domain.ErrUserNotFound
}

func (s *Store) FindByID(ctx context.Context, id string) (*domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if id == "" {
		return nil, domain.ErrInvalidInput
	}

	user, ok := s.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	return cloneUser(user), nil
}
