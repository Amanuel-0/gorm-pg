package models

import (
	"time"

	"gorm.io/gorm"
)

// create an enum for the message type

type MessageType string

const (
	MessageTypeText   MessageType = "text"
	MessageTypeImage  MessageType = "image"
	MessageTypeFile   MessageType = "file"
	MessageTypeSystem MessageType = "system"
)

func (mt MessageType) IsValid() bool {
	switch mt {
	case MessageTypeText, MessageTypeImage, MessageTypeFile, MessageTypeSystem:
		return true
	}
	return false
}

type Message struct {
	// gorm.Model
	ID          uint        `json:"id,omitempty" gorm:"primaryKey"`
	ThreadID    uint        `json:"thread_id,omitempty" gorm:"index"`
	SenderID    uint        `json:"sender_id,omitempty" gorm:"index"`
	Type        MessageType `json:"type,omitempty" gorm:"type:enum('text','image','file','system')"`
	Body        string      `json:"body,omitempty" gorm:"type:text"`
	Attachments string      `json:"attachments,omitempty" gorm:"type:json"` // JSON array of attachment URLs
	Deleted     bool        `json:"deleted,omitempty" gorm:"default:false"`

	// timestamps
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`

	// relationships
	Thread *ChatThread `json:"thread,omitempty" gorm:"foreignKey:ThreadID"`
	Sender *User       `json:"sender,omitempty" gorm:"foreignKey:SenderID"`
}
