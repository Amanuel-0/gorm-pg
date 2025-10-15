package models

import (
	"time"

	"gorm.io/gorm"
)

type Report struct {
	// gorm.Model
	ID         uint      `json:"id,omitempty" gorm:"primaryKey"`
	ReporterID uint      `json:"reporter_id,omitempty" gorm:"index"`
	TargetType string    `json:"target_type,omitempty" gorm:"size:50"`
	TargetID   uint      `json:"target_id,omitempty" gorm:"index"`
	Reason     string    `json:"reason,omitempty" gorm:"type:text"`
	Metadata   string    `json:"metadata,omitempty" gorm:"type:json"`
	HandledBy  uint      `json:"handled_by,omitempty" gorm:"index"`
	HandledAt  time.Time `json:"handled_at,omitempty" gorm:"autoUpdateTime"`
	Resolution string    `json:"resolution,omitempty" gorm:"type:text"`

	// timestamps
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`

	// Relationships
	Reporter User `json:"reporter,omitempty" gorm:"foreignKey:ReporterID"`
	Handler  User `json:"handler,omitempty" gorm:"foreignKey:HandledBy"`
}
