package models

import (
	"time"

	"gorm.io/gorm"
)

// create an enum for the condition
type Condition string

const (
	ConditionNew        Condition = "new"
	ConditionLikeNew    Condition = "like_new"
	ConditionGood       Condition = "good"
	ConditionAcceptable Condition = "acceptable"
)

func (c Condition) IsValid() bool {
	switch c {
	case ConditionNew, ConditionLikeNew, ConditionGood, ConditionAcceptable:
		return true
	}
	return false
}

type Book struct {
	gorm.Model
	OwnerID         uint       `json:"owner_id" gorm:"not null"`
	Title           string     `json:"title" gorm:"type:varchar(1000);not null"`
	Subtitle        *string    `json:"subtitle" gorm:"type:varchar(1000)"`
	AuthorID        *uint      `json:"author_id"`
	ISBN            *string    `json:"isbn" gorm:"type:varchar(32)"`
	Description     *string    `json:"description" gorm:"type:text"`
	Language        string     `json:"language" gorm:"type:varchar(8);default:'EN'"`
	Condition       string     `json:"condition" gorm:"type:enum('new','like_new','good','acceptable');default:'good'"`
	AvailableFrom   *time.Time `json:"available_from"`
	AvailableUntil  *time.Time `json:"available_until"`
	LocationCity    *string    `json:"location_city" gorm:"type:varchar(100)"`
	LocationState   *string    `json:"location_state" gorm:"type:varchar(100)"`
	LocationCountry *string    `json:"location_country" gorm:"type:varchar(100)"`
	Latitude        *float64   `json:"latitude"`
	Longitude       *float64   `json:"longitude"`
	Location        *string    `json:"location" gorm:"type:point"`
	Active          bool       `json:"active" gorm:"default:true"`
	ArchivedAt      *time.Time `json:"archived_at"`
	PreferredTitles *string    `json:"preferred_titles" gorm:"type:json"`

	// Relationships
	Owner  User        `json:"owner" gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE"`
	Author *Author     `json:"author" gorm:"foreignKey:AuthorID;constraint:OnDelete:SET NULL"`
	Genres []Genre     `json:"genres" gorm:"many2many:book_genres;"`
	Images []BookImage `json:"images" gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE"`
}
