package repository

import (
	"context"
	"errors"
	"sync"

	"supply-chain-aggregator/services/auth-service/internal/entity"
)

var ErrUserNotFound = errors.New("user not found")
var ErrEmailAlreadyExists = errors.New("email already exists")

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}

type InMemoryUserRepository struct {
	mu           sync.RWMutex
	usersByEmail map[string]*entity.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		usersByEmail: make(map[string]*entity.User),
	}
}

func (r *InMemoryUserRepository) Create(ctx context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if _, exists := r.usersByEmail[user.Email]; exists {
		return ErrEmailAlreadyExists
	}

	r.usersByEmail[user.Email] = user

	return nil
}

func (r *InMemoryUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	user, exists := r.usersByEmail[email]
	if !exists {
		return nil, ErrUserNotFound
	}

	return user, nil
}
