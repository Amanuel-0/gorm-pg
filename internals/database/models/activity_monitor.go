package models

import (
	"time"

	"gorm.io/gorm"
)

type ActivityAction string

const (
	LogActionCreate ActivityAction = "create"
	LogActionUpdate ActivityAction = "update"
	LogActionDelete ActivityAction = "delete"
	LogActionLogin  ActivityAction = "login"
	LogActionLogout ActivityAction = "logout"
)

func (a ActivityAction) IsValid() bool {
	switch a {
	case LogActionCreate, LogActionUpdate, LogActionDelete, LogActionLogin, LogActionLogout:
		return true
	}
	return false
}

type ActivityLog struct {
	// gorm.Model
	ID         uint           `json:"id,omitempty" gorm:"primaryKey"`
	UserID     *uint          `json:"user_id,omitempty" gorm:"index"`
	Action     ActivityAction `json:"action,omitempty" gorm:"type:enum('create','update','delete','login','logout');size:200;not null"`
	ObjectType string         `json:"object_type,omitempty" gorm:"size:100;index:idx_object"`
	ObjectID   *uint          `json:"object_id,omitempty" gorm:"index:idx_object"`
	Payload    string         `json:"payload,omitempty" gorm:"type:json"`

	IPAddress string `json:"ip_address,omitempty" gorm:"size:45"`
	UserAgent string `json:"user_agent,omitempty" gorm:"size:255"`
	RequestID string `json:"request_id,omitempty" gorm:"size:100"`

	// timestamps
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
