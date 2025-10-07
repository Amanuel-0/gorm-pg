package models

import (
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
	gorm.Model
	ThreadID    uint        `json:"thread_id" gorm:"index"`
	SenderID    uint        `json:"sender_id" gorm:"index"`
	Type        MessageType `json:"type" gorm:"type:enum('text','image','file','system')"`
	Body        string      `json:"body" gorm:"type:text"`
	Attachments string      `json:"attachments" gorm:"type:json"` // JSON array of attachment URLs
	Deleted     bool        `json:"deleted" gorm:"default:false"`

	// relationships
	Thread ChatThread `json:"thread" gorm:"foreignKey:ThreadID"`
	Sender User       `json:"sender" gorm:"foreignKey:SenderID"`
}
