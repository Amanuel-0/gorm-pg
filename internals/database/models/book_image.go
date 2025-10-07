package models

import (
	"time"

	"gorm.io/gorm"
)

type BookImage struct {
	gorm.Model
	URL        string    `json:"url"`
	Width      int       `json:"width"`
	Height     int       `json:"height"`
	IsPrimary  bool      `json:"is_primary"`
	UploadedAt time.Time `json:"uploaded_at"`

	// relationships
	BookID uint `json:"book_id"`
	Book   Book `json:"book"`
}
