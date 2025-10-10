package models

import (
	"time"

	"gorm.io/gorm"
)

type Genre struct {
	// gorm.Model
	ID        uint            `json:"id" gorm:"primaryKey"`
	Slug      string          `json:"slug,omitempty"`
	Name      string          `json:"name,omitempty"`
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`
}
