package models

import (
	"time"

	"gorm.io/gorm"
)

type CommunityMessage struct {
	// gorm.Model
	ID       uint   `json:"id,omitempty" gorm:"primaryKey"`
	ThreadID uint   `json:"thread_id,omitempty" gorm:"index;not null"`
	SenderID uint   `json:"sender_id,omitempty" gorm:"index;not null"`
	Body     string `json:"body,omitempty" gorm:"type:text;not null"`

	// timestamps
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`

	// Relationships
	Thread *CommunityThread `json:"thread,omitempty" gorm:"foreignKey:ThreadID"`
	Sender *User            `json:"sender,omitempty" gorm:"foreignKey:SenderID"`
}
