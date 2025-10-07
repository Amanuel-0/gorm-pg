package models

import (
	"gorm.io/gorm"
)

type ActivityAction string

const (
	ActionCreate ActivityAction = "create"
	ActionUpdate ActivityAction = "update"
	ActionDelete ActivityAction = "delete"
	ActionLogin  ActivityAction = "login"
	ActionLogout ActivityAction = "logout"
)

func (a ActivityAction) IsValid() bool {
	switch a {
	case ActionCreate, ActionUpdate, ActionDelete, ActionLogin, ActionLogout:
		return true
	}
	return false
}

type ActivityLog struct {
	gorm.Model
	UserID     *uint          `json:"user_id" gorm:"index"`
	Action     ActivityAction `json:"action" gorm:"type:enum('create','update','delete','login','logout');size:200;not null"`
	ObjectType string         `json:"object_type" gorm:"size:100;index:idx_object"`
	ObjectID   *uint          `json:"object_id" gorm:"index:idx_object"`
	Payload    string         `json:"payload" gorm:"type:json"`

	IPAddress string `json:"ip_address" gorm:"size:45"`
	UserAgent string `json:"user_agent" gorm:"size:255"`
	RequestID string `json:"request_id" gorm:"size:100"`

	// Relationships
	User *User `json:"user" gorm:"foreignKey:UserID"`
}
