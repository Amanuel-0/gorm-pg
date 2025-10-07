package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

/** enums starts */
type Interval string

const (
	IntervalMonth  Interval = "month"
	Interval3Month Interval = "3_month"
	IntervalYear   Interval = "year"
)

func (i Interval) IsValid() bool {
	switch i {
	case IntervalMonth, Interval3Month, IntervalYear:
		return true
	}
	return false
}

/** enums ends */

type SubscriptionPlan struct {
	gorm.Model
	Slug        string         `gorm:"size:100;not null;unique"`
	Name        string         `gorm:"size:200;not null"`
	Description string         `gorm:"type:text"`
	PriceCents  int            `gorm:"not null"`
	Currency    string         `gorm:"size:8;not null;default:USD"`
	Interval    string         `gorm:"type:enum('month','year');not null"`
	Features    datatypes.JSON `gorm:"type:json"`
	Active      bool           `gorm:"type:tinyint(1);not null;default:1"`
}
