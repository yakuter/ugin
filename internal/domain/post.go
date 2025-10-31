package domain

import "time"

// Post represents a blog post or article
type Post struct {
	ID          uint       `json:"id" gorm:"primarykey" example:"1"`
	CreatedAt   time.Time  `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time  `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null" example:"Getting Started with Go"`
	Description string     `json:"description" gorm:"type:text" example:"A comprehensive guide to learning Go programming language"`
	Tags        []Tag      `json:"tags,omitempty" gorm:"foreignKey:PostID"`
}

// Tag represents a tag associated with a post
type Tag struct {
	ID          uint       `json:"id" gorm:"primarykey" example:"1"`
	CreatedAt   time.Time  `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time  `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
	PostID      uint       `json:"post_id" gorm:"index;not null" example:"1"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null" example:"golang"`
	Description string     `json:"description" gorm:"type:text" example:"Go programming language"`
}

// TableName overrides the table name for Post
func (Post) TableName() string {
	return "posts"
}

// TableName overrides the table name for Tag
func (Tag) TableName() string {
	return "tags"
}

// CreatePostRequest represents the request body for creating a post
type CreatePostRequest struct {
	Name        string           `json:"name" binding:"required" example:"Getting Started with Go"`
	Description string           `json:"description" example:"A comprehensive guide to learning Go programming language"`
	Tags        []CreateTagRequest `json:"tags,omitempty"`
}

// CreateTagRequest represents a tag in the create request
type CreateTagRequest struct {
	Name        string `json:"name" binding:"required" example:"golang"`
	Description string `json:"description" example:"Go programming language"`
}

