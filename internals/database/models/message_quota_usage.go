package models

import (
	"time"

	"gorm.io/gorm"
)

type MessageQuotaUsage struct {
	gorm.Model
	UserID       uint      `json:"user_id" gorm:"index"`
	PeriodStart  time.Time `json:"period_start" gorm:"type:date"`
	PeriodEnd    time.Time `json:"period_end" gorm:"type:date"`
	MessagesSent int       `json:"messages_sent" gorm:"default:0"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
