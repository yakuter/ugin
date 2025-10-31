package domain

import "time"

// Post represents a blog post or article
type Post struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text"`
	Tags        []Tag     `json:"tags,omitempty" gorm:"foreignKey:PostID"`
}

// Tag represents a tag associated with a post
type Tag struct {
	ID          uint       `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	PostID      uint       `json:"post_id" gorm:"index;not null"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null"`
	Description string     `json:"description" gorm:"type:text"`
}

// TableName overrides the table name for Post
func (Post) TableName() string {
	return "posts"
}

// TableName overrides the table name for Tag
func (Tag) TableName() string {
	return "tags"
}

