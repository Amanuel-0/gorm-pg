package models

import (
	"gorm.io/gorm"
)

type Genre struct {
	gorm.Model
	Slug string `json:"slug"`
	Name string `json:"name"`
}
