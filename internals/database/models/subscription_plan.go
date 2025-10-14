package models

import (
	"time"

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
	// gorm.Model
	ID          uint           `json:"id,omitempty" gorm:"primaryKey"`
	Slug        string         `json:"slug,omitempty" gorm:"size:100;not null;unique"`
	Name        string         `json:"name,omitempty" gorm:"size:200;not null"`
	Description string         `json:"description,omitempty" gorm:"type:text"`
	PriceCents  int            `json:"price_cents,omitempty" gorm:"not null"`
	Currency    string         `json:"currency,omitempty" gorm:"size:8;not null;default:USD"`
	Interval    Interval       `json:"interval,omitempty" gorm:"type:enum('month','3_month','year');not null"`
	Features    datatypes.JSON `json:"features,omitempty" gorm:"type:json"`
	Active      bool           `json:"active,omitempty" gorm:"type:tinyint(1);not null;default:1"`

	// timestamp
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`
}
