package models

import (
	"gorm.io/gorm"
)

type CommunityMessage struct {
	gorm.Model
	ThreadID uint   `json:"thread_id" gorm:"index;not null"`
	SenderID uint   `json:"sender_id" gorm:"index;not null"`
	Body     string `json:"body" gorm:"type:text;not null"`

	// Relationships
	Thread *CommunityThread `json:"thread" gorm:"foreignKey:ThreadID"`
	Sender *User            `json:"sender" gorm:"foreignKey:SenderID"`
}
