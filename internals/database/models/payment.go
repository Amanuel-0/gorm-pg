package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSucceeded PaymentStatus = "succeeded"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusCanceled  PaymentStatus = "canceled"
)

type Payment struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	UserID         uint           `json:"user_id"`
	SubscriptionID uint           `json:"subscription_id"`
	AmountCents    float64        `json:"amount_cents"`
	Status         PaymentStatus  `json:"status" gorm:"type:enum('pending','succeeded','failed','refunded','canceled');not null;default:'pending'"`
	Metadata       datatypes.JSON `json:"metadata" gorm:"type:json"`

	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`

	// relationships
	User         User         `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Subscription Subscription `json:"subscription" gorm:"foreignKey:SubscriptionID;constraint:OnDelete:CASCADE"`
}
