package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yakuter/ugin/internal/domain"
	"github.com/yakuter/ugin/internal/repository"
	"github.com/yakuter/ugin/internal/service"
)

// Mock repository
type mockPostRepository struct {
	getByIDFunc func(ctx context.Context, id string) (*domain.Post, error)
	listFunc    func(ctx context.Context, filter repository.ListFilter) ([]*domain.Post, *repository.ListResult, error)
	createFunc  func(ctx context.Context, post *domain.Post) error
	updateFunc  func(ctx context.Context, post *domain.Post) error
	deleteFunc  func(ctx context.Context, id string) error
}

func (m *mockPostRepository) GetByID(ctx context.Context, id string) (*domain.Post, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockPostRepository) List(ctx context.Context, filter repository.ListFilter) ([]*domain.Post, *repository.ListResult, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, filter)
	}
	return nil, nil, errors.New("not implemented")
}

func (m *mockPostRepository) Create(ctx context.Context, post *domain.Post) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, post)
	}
	return errors.New("not implemented")
}

func (m *mockPostRepository) Update(ctx context.Context, post *domain.Post) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, post)
	}
	return errors.New("not implemented")
}

func (m *mockPostRepository) Delete(ctx context.Context, id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

// Mock logger
type mockLogger struct{}

func (m *mockLogger) Debug(msg string, keysAndValues ...interface{}) {}
func (m *mockLogger) Info(msg string, keysAndValues ...interface{})  {}
func (m *mockLogger) Warn(msg string, keysAndValues ...interface{})  {}
func (m *mockLogger) Error(msg string, keysAndValues ...interface{}) {}

func TestPostService_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mock    func() *mockPostRepository
		wantErr bool
	}{
		{
			name: "success",
			id:   "1",
			mock: func() *mockPostRepository {
				return &mockPostRepository{
					getByIDFunc: func(ctx context.Context, id string) (*domain.Post, error) {
						return &domain.Post{
							ID:          1,
							Name:        "Test Post",
							Description: "Test Description",
						}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "not found",
			id:   "999",
			mock: func() *mockPostRepository {
				return &mockPostRepository{
					getByIDFunc: func(ctx context.Context, id string) (*domain.Post, error) {
						return nil, repository.ErrNotFound
					},
				}
			},
			wantErr: true,
		},
		{
			name:    "empty id",
			id:      "",
			mock:    func() *mockPostRepository { return &mockPostRepository{} },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewPostService(tt.mock(), &mockLogger{})
			post, err := svc.GetByID(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if post == nil {
				t.Error("expected post, got nil")
			}
		})
	}
}

func TestPostService_Create(t *testing.T) {
	tests := []struct {
		name    string
		post    *domain.Post
		mock    func() *mockPostRepository
		wantErr bool
	}{
		{
			name: "success",
			post: &domain.Post{
				Name:        "New Post",
				Description: "Description",
			},
			mock: func() *mockPostRepository {
				return &mockPostRepository{
					createFunc: func(ctx context.Context, post *domain.Post) error {
						post.ID = 1
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:    "nil post",
			post:    nil,
			mock:    func() *mockPostRepository { return &mockPostRepository{} },
			wantErr: true,
		},
		{
			name: "empty name",
			post: &domain.Post{
				Description: "Description without name",
			},
			mock:    func() *mockPostRepository { return &mockPostRepository{} },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewPostService(tt.mock(), &mockLogger{})
			err := svc.Create(context.Background(), tt.post)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

