package models

import (
	"gorm.io/gorm"
)

type UserProfile struct {
	gorm.Model
	ID          uint   `json:"id" gorm:"primaryKey"`
	FirstName   string `json:"first_name" gorm:"size:100"`
	LastName    string `json:"last_name" gorm:"size:100"`
	DisplayName string `json:"display_name" gorm:"size:25"`
	Bio         string `json:"bio" gorm:"size:255"`
	AvatarURL   string `json:"avatar_url" gorm:"size:1000"`

	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Location  string  `json:"location" gorm:"type:geometry"`

	Linkedin string `json:"linkedin"`

	// relationships
	UserID    uint     `json:"user_id" gorm:"uniqueIndex"`
	CountryID *uint    `json:"country_id"`
	Country   *Country `json:"country" gorm:"foreignKey:CountryID"`
	StateID   *uint    `json:"state_id"`
	State     *State   `json:"state" gorm:"foreignKey:StateID"`
	CityID    *uint    `json:"city_id"`
	City      *State   `json:"city" gorm:"foreignKey:CityID"`
}
