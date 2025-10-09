package models

import (
	"time"

	"gorm.io/gorm"
)

type UserProfile struct {
	// gorm.Model
	ID          uint   `json:"id,omitempty" gorm:"primaryKey"`
	FirstName   string `json:"first_name,omitempty" gorm:"size:100"`
	LastName    string `json:"last_name,omitempty" gorm:"size:100"`
	DisplayName string `json:"display_name,omitempty" gorm:"size:25"`
	Bio         string `json:"bio,omitempty" gorm:"size:255"`
	AvatarURL   string `json:"avatar_url,omitempty" gorm:"size:1000"`

	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Location  string  `json:"location,omitempty" gorm:"type:geometry"`

	Linkedin string `json:"linkedin,omitempty"`

	// timestamps
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`

	// relationships
	UserID    uint     `json:"user_id,omitempty" gorm:"uniqueIndex"`
	CountryID *uint    `json:"country_id,omitempty"`
	Country   *Country `json:"country,omitempty" gorm:"foreignKey:CountryID"`
	StateID   *uint    `json:"state_id,omitempty"`
	State     *State   `json:"state,omitempty" gorm:"foreignKey:StateID"`
	CityID    *uint    `json:"city_id,omitempty"`
	City      *State   `json:"city,omitempty" gorm:"foreignKey:CityID"`
}
