package models

import (
	"gorm.io/gorm"
)

type UserRating struct {
	gorm.Model
	RaterID     uint   `json:"raterId" gorm:"index"`
	RatedUserID uint   `json:"ratedUserId" gorm:"index"`
	ExchangeID  uint   `json:"exchangeId" gorm:"index"`
	Rating      uint8  `json:"rating" gorm:"check:rating >= 1 AND rating <= 5"`
	Comment     string `json:"comment" gorm:"type:text"`

	// Relationships
	Rater     User     `json:"rater" gorm:"foreignKey:RaterID"`
	RatedUser User     `json:"ratedUser" gorm:"foreignKey:RatedUserID"`
	Exchange  Exchange `json:"exchange" gorm:"foreignKey:ExchangeID"`
}
