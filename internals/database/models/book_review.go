package models

import (
	"gorm.io/gorm"
)

type BookReview struct {
	gorm.Model
	BookID     uint   `json:"book_id" gorm:"index;not null"`
	ReviewerID uint   `json:"reviewer_id" gorm:"index;not null"`
	Rating     uint   `json:"rating" gorm:"not null;check:rating >=1 AND rating <=5"`
	Comment    string `json:"comment" gorm:"type:text;null"`

	// Relationships
	Book Book `json:"book" gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE"`
	User User `json:"user" gorm:"foreignKey:ReviewerID;constraint:OnDelete:CASCADE"`
}
