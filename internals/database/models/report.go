package models

import (
	"time"

	"gorm.io/gorm"
)

type Report struct {
	gorm.Model
	ReporterID uint      `json:"reporter_id" gorm:"index"`
	TargetType string    `json:"target_type" gorm:"size:50"`
	TargetID   uint      `json:"target_id" gorm:"index"`
	Reason     string    `json:"reason" gorm:"type:text"`
	Metadata   string    `json:"metadata" gorm:"type:json"`
	HandledBy  uint      `json:"handled_by" gorm:"index"`
	HandledAt  time.Time `json:"handled_at" gorm:"autoUpdateTime"`
	Resolution string    `json:"resolution" gorm:"type:text"`

	// Relationships
	Reporter User `json:"reporter,omitempty" gorm:"foreignKey:ReporterID"`
	Handler  User `json:"handler,omitempty" gorm:"foreignKey:HandledBy"`
}
