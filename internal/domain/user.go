package domain

import "time"

// User represents a user in the system
type User struct {
	ID             uint       `json:"id" gorm:"primarykey"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	Email          string     `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	MasterPassword string     `json:"-" gorm:"type:varchar(255);not null"` // Never expose in JSON
}

// TableName overrides the table name for User
func (User) TableName() string {
	return "users"
}

