package models

import (
	"time"

	"gorm.io/gorm"
)

type CommunityRole string

const (
	CommunityRoleMember    CommunityRole = "member"
	CommunityRoleAdmin     CommunityRole = "admin"
	CommunityRoleModerator CommunityRole = "moderator"
)

func (cr CommunityRole) IsValid() bool {
	switch cr {
	case CommunityRoleMember, CommunityRoleAdmin, CommunityRoleModerator:
		return true
	}
	return false
}

type CommunityMember struct {
	gorm.Model
	// Composite unique index on CommunityID and UserID
	CommunityID   uint          `json:"community_id" gorm:"uniqueIndex:idx_community_user;not null"`
	UserID        uint          `json:"user_id" gorm:"uniqueIndex:idx_community_user;not null"`
	CommunityRole CommunityRole `json:"community_role" gorm:"type:enum('member','admin','moderator');default:'member'"`
	JoinedAt      time.Time     `json:"joined_at" gorm:"autoCreateTime"`

	// Relationships
	Community *Community `json:"community,omitempty" gorm:"foreignKey:CommunityID"`
	User      *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
