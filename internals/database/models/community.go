package models

import (
	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name            string `json:"name" gorm:"size:255;not null"`
	Slug            string `json:"slug" gorm:"size:255;unique"`
	Description     string `json:"description" gorm:"type:text;not null"`
	CreatorID       uint   `json:"creator_id" gorm:"index"`
	RequirePaidChat bool   `json:"require_paid_chat" gorm:"default:true"`

	// Relationships
	Creator User              `json:"creator" gorm:"foreignKey:CreatorID"`
	Members []CommunityMember `json:"members" gorm:"foreignKey:CommunityID"`
	Threads []CommunityThread `json:"threads" gorm:"foreignKey:CommunityID"`
}
