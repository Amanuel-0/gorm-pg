package models

import (
	"gorm.io/gorm"
)

type Country struct {
	gorm.Model
	Name string `json:"name"`
	Code string `json:"code"`

	// One-to-Many => Country has many States
	States []State `json:"states" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type State struct {
	gorm.Model
	Name string `json:"name"`

	// Many-to-One => State belongs to Country
	CountryID uint    `json:"country_id"`
	Country   Country `json:"country"`

	// One-to-Many → State has many Cities
	Cities []City `json:"cities" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type City struct {
	gorm.Model
	Name string `json:"name"`

	// Many-to-One → City belongs to State
	StateID uint  `json:"state_id"`
	State   State `json:"state"`
}
