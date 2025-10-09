package models

import (
	"time"

	"gorm.io/gorm"
)

type Author struct {
	// gorm.Model
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name,omitempty" gorm:"uniqueIndex"`

	// timestamps
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`
}
