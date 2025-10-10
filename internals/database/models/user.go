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
	// gorm.Model
	ID              uint       `json:"id,omitempty" gorm:"primaryKey"`
	Email           string     `json:"email,omitempty" gorm:"unique;not null"`
	Phone           string     `json:"phone,omitempty" gorm:"unique;not null"`
	PasswordHash    string     `json:"-" gorm:"column:password_hash;not null"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty" gorm:"column:email_verified_at"`
	PhoneVerifiedAt *time.Time `json:"phone_verified_at,omitempty" gorm:"column:phone_verified_at"`
	FirstName       string     `json:"first_name,omitempty" gorm:"size:100"`
	LastName        string     `json:"last_name,omitempty" gorm:"size:100"`
	IsActive        bool       `json:"is_active,omitempty" gorm:"default:true"`
	Role            Role       `json:"role,omitempty" gorm:"type:enum('user','admin','moderator','system');default:'user'"`
	Local           string     `json:"local,omitempty" gorm:"default:'en'"`

	// timestamps
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`

	// Relationships
	// UserProfile UserProfile `json:"user_profile" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserProfile UserProfile `json:"user_profile,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// user_preferred_genres which is a many to many relationship
	PreferredGenres []*Genre `json:"preferred_genres,omitempty" gorm:"many2many:user_preferred_genres"`
}
