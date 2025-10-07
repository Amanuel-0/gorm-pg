package models

import (
	"gorm.io/gorm"
)

type ChatThread struct {
	gorm.Model
	ExchangeID uint `json:"exchange_id" gorm:"index"`
	CreatedBy  uint `json:"created_by" gorm:"index"`
	Archived   bool `json:"archived" gorm:"default:false"`

	// relationships
	Exchange Exchange  `json:"exchange" gorm:"foreignKey:ExchangeID"`
	Creator  User      `json:"creator" gorm:"foreignKey:CreatedBy"`
	Messages []Message `json:"messages" gorm:"foreignKey:ThreadID;constraint:OnDelete:CASCADE"`
}
