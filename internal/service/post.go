package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/yakuter/ugin/internal/domain"
	"github.com/yakuter/ugin/internal/repository"
)

type postService struct {
	repo   repository.PostRepository
	logger Logger
}

// NewPostService creates a new post service
func NewPostService(repo repository.PostRepository, logger Logger) PostService {
	return &postService{
		repo:   repo,
		logger: logger,
	}
}

func (s *postService) GetByID(ctx context.Context, id string) (*domain.Post, error) {
	if id == "" {
		return nil, repository.ErrInvalidInput
	}

	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			s.logger.Info("post not found", "id", id)
			return nil, err
		}
		s.logger.Error("failed to get post", "id", id, "error", err)
		return nil, fmt.Errorf("get post: %w", err)
	}

	return post, nil
}

func (s *postService) List(ctx context.Context, filter repository.ListFilter) ([]*domain.Post, *repository.ListResult, error) {
	// Set default values
	if filter.Limit <= 0 {
		filter.Limit = 25
	}
	if filter.Limit > 100 {
		filter.Limit = 100 // Max limit
	}
	if filter.Sort == "" {
		filter.Sort = "id"
	}
	if filter.Order == "" {
		filter.Order = "DESC"
	}

	posts, result, err := s.repo.List(ctx, filter)
	if err != nil {
		s.logger.Error("failed to list posts", "error", err)
		return nil, nil, fmt.Errorf("list posts: %w", err)
	}

	s.logger.Debug("listed posts", "count", len(posts), "total", result.Total)
	return posts, result, nil
}

func (s *postService) Create(ctx context.Context, post *domain.Post) error {
	if post == nil {
		return repository.ErrInvalidInput
	}

	if post.Name == "" {
		return fmt.Errorf("%w: name is required", repository.ErrInvalidInput)
	}

	if err := s.repo.Create(ctx, post); err != nil {
		s.logger.Error("failed to create post", "error", err)
		return fmt.Errorf("create post: %w", err)
	}

	s.logger.Info("post created", "id", post.ID, "name", post.Name)
	return nil
}

func (s *postService) Update(ctx context.Context, id string, post *domain.Post) error {
	if id == "" || post == nil {
		return repository.ErrInvalidInput
	}

	// Check if post exists
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Update fields
	existing.Name = post.Name
	existing.Description = post.Description
	existing.Tags = post.Tags

	if err := s.repo.Update(ctx, existing); err != nil {
		s.logger.Error("failed to update post", "id", id, "error", err)
		return fmt.Errorf("update post: %w", err)
	}

	s.logger.Info("post updated", "id", id)
	return nil
}

func (s *postService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return repository.ErrInvalidInput
	}

	// Check if post exists
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete post", "id", id, "error", err)
		return fmt.Errorf("delete post: %w", err)
	}

	s.logger.Info("post deleted", "id", id)
	return nil
}

