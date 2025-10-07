package models

import (
	"gorm.io/gorm"
)

type CommunityThread struct {
	gorm.Model
	CommunityID uint   `json:"community_id" gorm:"index;not null"`
	CreatedBy   uint   `json:"created_by" gorm:"index"`
	Title       string `json:"title" gorm:"size:500"`

	Community *Community         `json:"community" gorm:"foreignKey:CommunityID"`
	Creator   *User              `json:"creator" gorm:"foreignKey:CreatedBy"`
	Messages  []CommunityMessage `json:"messages" gorm:"foreignKey:ThreadID"`
}
