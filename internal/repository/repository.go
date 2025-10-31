package repository

import (
	"context"
	"errors"

	"github.com/yakuter/ugin/internal/domain"
)

// Common errors
var (
	ErrNotFound      = errors.New("record not found")
	ErrAlreadyExists = errors.New("record already exists")
	ErrInvalidInput  = errors.New("invalid input")
)

// ListFilter contains common filtering options
type ListFilter struct {
	Search string
	Limit  int
	Offset int
	Sort   string
	Order  string
}

// ListResult contains paginated results
type ListResult struct {
	Total    int64
	Filtered int64
}

// PostRepository defines the interface for post data access
type PostRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Post, error)
	List(ctx context.Context, filter ListFilter) ([]*domain.Post, *ListResult, error)
	Create(ctx context.Context, post *domain.Post) error
	Update(ctx context.Context, post *domain.Post) error
	Delete(ctx context.Context, id string) error
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uint) error
}

