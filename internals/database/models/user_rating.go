package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRating struct {
	// gorm.Model
	ID          uint   `json:"id,omitempty" gorm:"primaryKey"`
	RaterID     uint   `json:"rater_id,omitempty" gorm:"index"`
	RatedUserID uint   `json:"rated_user_id,omitempty" gorm:"index"`
	ExchangeID  uint   `json:"exchange_id,omitempty" gorm:"index"`
	Rating      uint8  `json:"rating,omitempty" gorm:"check:rating >= 1 AND rating <= 5"`
	Comment     string `json:"comment,omitempty" gorm:"type:text"`

	// timestamps
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`

	// Relationships
	Rater     User     `json:"rater,omitempty" gorm:"foreignKey:RaterID"`
	RatedUser User     `json:"rated_user,omitempty" gorm:"foreignKey:RatedUserID"`
	Exchange  Exchange `json:"exchange,omitempty" gorm:"foreignKey:ExchangeID"`
}
