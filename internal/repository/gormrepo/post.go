package gormrepo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/yakuter/ugin/internal/domain"
	"github.com/yakuter/ugin/internal/repository"
	"gorm.io/gorm"
)

type postRepository struct {
	db *gorm.DB
}

// NewPostRepository creates a new post repository
func NewPostRepository(db *gorm.DB) repository.PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) GetByID(ctx context.Context, id string) (*domain.Post, error) {
	var post domain.Post

	err := r.db.WithContext(ctx).
		Preload("Tags").
		First(&post, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return &post, nil
}

func (r *postRepository) List(ctx context.Context, filter repository.ListFilter) ([]*domain.Post, *repository.ListResult, error) {
	var posts []*domain.Post
	result := &repository.ListResult{}

	query := r.db.WithContext(ctx).Model(&domain.Post{})

	// Apply search filter
	if filter.Search != "" {
		searchTerm := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchTerm, searchTerm)
	}

	// Get filtered count
	if err := query.Count(&result.Filtered).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to count filtered posts: %w", err)
	}

	// Get total count (without filters)
	if err := r.db.WithContext(ctx).Model(&domain.Post{}).Count(&result.Total).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to count total posts: %w", err)
	}

	// Apply sorting
	if filter.Sort != "" {
		order := "ASC"
		if strings.ToUpper(filter.Order) == "DESC" {
			order = "DESC"
		}
		// Sanitize sort field to prevent SQL injection
		sortField := strings.ToLower(filter.Sort)
		if sortField == "id" || sortField == "name" || sortField == "created_at" || sortField == "updated_at" {
			query = query.Order(fmt.Sprintf("%s %s", sortField, order))
		}
	}

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	// Fetch results with tags
	if err := query.Preload("Tags").Find(&posts).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to list posts: %w", err)
	}

	return posts, result, nil
}

func (r *postRepository) Create(ctx context.Context, post *domain.Post) error {
	if err := r.db.WithContext(ctx).Create(post).Error; err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}
	return nil
}

func (r *postRepository) Update(ctx context.Context, post *domain.Post) error {
	if err := r.db.WithContext(ctx).Save(post).Error; err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	return nil
}

func (r *postRepository) Delete(ctx context.Context, id string) error {
	// Start a transaction to delete post and its tags
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Soft delete tags associated with the post
		if err := tx.Where("post_id = ?", id).Delete(&domain.Tag{}).Error; err != nil {
			return fmt.Errorf("failed to delete tags: %w", err)
		}

		// Soft delete the post
		if err := tx.Delete(&domain.Post{}, id).Error; err != nil {
			return fmt.Errorf("failed to delete post: %w", err)
		}

		return nil
	})
}

