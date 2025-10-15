package models

import (
	"time"

	"gorm.io/gorm"
)

type Community struct {
	// gorm.Model
	ID              uint   `json:"id,omitempty" gorm:"primaryKey"`
	Name            string `json:"name,omitempty" gorm:"size:255;not null"`
	Slug            string `json:"slug,omitempty" gorm:"size:255;unique"`
	Description     string `json:"description,omitempty" gorm:"type:text;not null"`
	CreatorID       uint   `json:"creator_id,omitempty" gorm:"index"`
	RequirePaidChat bool   `json:"require_paid_chat,omitempty" gorm:"default:true"`

	// timestamps
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`

	// Relationships
	Creator User              `json:"creator,omitempty" gorm:"foreignKey:CreatorID"`
	Members []CommunityMember `json:"members,omitempty" gorm:"foreignKey:CommunityID"`
	Threads []CommunityThread `json:"threads,omitempty" gorm:"foreignKey:CommunityID"`
}
