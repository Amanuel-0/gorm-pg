package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleUser      Role = "user"
	RoleAdmin     Role = "admin"
	RoleModerator Role = "moderator"
	RoleSystem    Role = "system"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleUser, RoleAdmin, RoleModerator, RoleSystem:
		return true
	}
	return false
}

type User struct {
	gorm.Model
	ID              uint       `json:"id" gorm:"primaryKey"`
	Email           string     `json:"email" gorm:"unique;not null"`
	Phone           string     `json:"phone" gorm:"unique;not null"`
	PasswordHash    string     `json:"-" gorm:"column:password_hash;not null"`
	EmailVerifiedAt *time.Time `json:"email_verified_at" gorm:"column:email_verified_at"`
	PhoneVerifiedAt *time.Time `json:"phone_verified_at" gorm:"column:phone_verified_at"`
	FirstName       string     `json:"first_name" gorm:"size:100"`
	LastName        string     `json:"last_name" gorm:"size:100"`
	IsActive        bool       `json:"is_active" gorm:"default:true"`
	Role            Role       `json:"role" gorm:"type:enum('user','admin','moderator','system');default:'user'"`
	Local           string     `json:"local" gorm:"default:'en'"`

	// Relationships
	// UserProfile UserProfile `json:"user_profile" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserProfile UserProfile `json:"user_profile" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
