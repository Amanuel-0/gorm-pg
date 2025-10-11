package models

import (
	"time"

	"gorm.io/gorm"
)

type ChatThread struct {
	// gorm.Model
	ID         uint `json:"id,omitempty" gorm:"primaryKey"`
	ExchangeID uint `json:"exchange_id,omitempty" gorm:"index"`
	CreatedBy  uint `json:"created_by,omitempty" gorm:"index"`
	Archived   bool `json:"archived,omitempty" gorm:"default:false"`

	// timestamps
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`

	// relationships
	Exchange *Exchange  `json:"exchange,omitempty" gorm:"foreignKey:ExchangeID"`
	Creator  *User      `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	Messages []*Message `json:"messages,omitempty" gorm:"foreignKey:ThreadID;constraint:OnDelete:CASCADE"`
}
