package models

import (
	"time"

	"gorm.io/gorm"
)

type CommunityThread struct {
	// gorm.Model
	ID          uint   `json:"id" gorm:"primaryKey"`
	CommunityID uint   `json:"community_id,omitempty" gorm:"index;not null"`
	CreatedBy   uint   `json:"created_by,omitempty" gorm:"index"`
	Title       string `json:"title,omitempty" gorm:"size:500"`

	// timestamps
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`

	// relationship
	Community *Community         `json:"community,omitempty" gorm:"foreignKey:CommunityID"`
	Creator   *User              `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	Messages  []CommunityMessage `json:"messages,omitempty" gorm:"foreignKey:ThreadID"`
	// Messages  []CommunityMessage `json:"messages,omitempty" gorm:"foreignKey:ThreadID"`
}
