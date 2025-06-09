package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user from Clerk authentication
type User struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ClerkUserID     string         `json:"clerkUserId" gorm:"uniqueIndex;not null"` // Clerk's user ID
	Email           string         `json:"email" gorm:"not null"`
	FirstName       string         `json:"firstName"`
	LastName        string         `json:"lastName"`
	ProfileImageURL string         `json:"profileImageUrl"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Schemas []Schema `json:"schemas,omitempty" gorm:"foreignKey:UserID"`
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return u.Email
	}
	return u.FirstName + " " + u.LastName
} 