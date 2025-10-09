package models

import (
	"time"

	"gorm.io/gorm"
)

type BookImage struct {
	// gorm.Model
	ID         uint      `json:"id,omitempty" gorm:"primaryKey"`
	URL        string    `json:"url"`
	Width      int       `json:"width"`
	Height     int       `json:"height"`
	IsPrimary  bool      `json:"is_primary"`
	UploadedAt time.Time `json:"uploaded_at"`

	// timestamps
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`

	// relationships
	BookID uint `json:"book_id"`
	Book   Book `json:"book"`
}
