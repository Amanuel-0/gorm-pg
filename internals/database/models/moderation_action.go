package models

import (
	"time"

	"gorm.io/gorm"
)

type ModerationAction struct {
	gorm.Model
	TargetType  string    `json:"target_type" gorm:"size:50"`
	TargetID    uint      `json:"target_id" gorm:"index"`
	Action      string    `json:"action" gorm:"size:100"`
	PerformedBy *uint     `json:"performed_by" gorm:"index"`
	PerformedAt time.Time `json:"performed_at" gorm:"autoCreateTime"`
	Reason      string    `json:"reason" gorm:"type:text"`
	Metadata    string    `json:"metadata" gorm:"type:json"`

	// Relationships
	Performer *User `gorm:"foreignKey:PerformedBy"`
}
