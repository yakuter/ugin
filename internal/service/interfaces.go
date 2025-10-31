package service

import (
	"context"

	"github.com/yakuter/ugin/internal/domain"
	"github.com/yakuter/ugin/internal/repository"
)

// Logger defines the logging interface
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

// PostService defines the business logic for posts
type PostService interface {
	GetByID(ctx context.Context, id string) (*domain.Post, error)
	List(ctx context.Context, filter repository.ListFilter) ([]*domain.Post, *repository.ListResult, error)
	Create(ctx context.Context, post *domain.Post) error
	Update(ctx context.Context, id string, post *domain.Post) error
	Delete(ctx context.Context, id string) error
}

// AuthService defines the business logic for authentication
type AuthService interface {
	SignIn(ctx context.Context, creds *domain.Credentials) (*domain.TokenDetails, error)
	SignUp(ctx context.Context, creds *domain.Credentials) error
	RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenDetails, error)
	ValidateToken(ctx context.Context, token string) (*domain.TokenClaims, error)
}

