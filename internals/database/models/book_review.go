package models

import (
	"time"

	"gorm.io/gorm"
)

type BookReview struct {
	// gorm.Model
	ID         uint   `json:"id,omitempty" gorm:"primaryKey"`
	BookID     uint   `json:"book_id,omitempty" gorm:"index;not null"`
	ReviewerID uint   `json:"reviewer_id,omitempty" gorm:"index;not null"`
	Rating     uint   `json:"rating,omitempty" gorm:"not null;check:rating >=1 AND rating <=5"`
	Comment    string `json:"comment,omitempty" gorm:"type:text;null"`

	// timestamps
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`

	// Relationships
	Book *Book `json:"book,omitempty" gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE"`
	User *User `json:"user,omitempty" gorm:"foreignKey:ReviewerID;constraint:OnDelete:CASCADE"`
}
