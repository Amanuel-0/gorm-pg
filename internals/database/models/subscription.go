package models

import (
	"time"

	"gorm.io/gorm"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive   SubscriptionStatus = "active"
	SubscriptionStatusPastDue  SubscriptionStatus = "past_due"
	SubscriptionStatusCanceled SubscriptionStatus = "canceled"
	SubscriptionStatusTrialing SubscriptionStatus = "trialing"
	SubscriptionStatusExpired  SubscriptionStatus = "expired"
)

func (s SubscriptionStatus) IsValid() bool {
	switch s {
	case SubscriptionStatusActive, SubscriptionStatusPastDue, SubscriptionStatusCanceled, SubscriptionStatusTrialing, SubscriptionStatusExpired:
		return true
	}
	return false
}

type Subscription struct {
	gorm.Model
	UserID                 uint64             `gorm:"not null;index"`
	PlanID                 *uint64            `gorm:"index"`
	ProviderSubscriptionID string             `gorm:"size:255;index"`
	Status                 SubscriptionStatus `gorm:"type:enum('active','past_due','canceled','trialing','expired');not null;default:'trialing'"`
	CurrentPeriodStart     *time.Time         `gorm:"type:datetime"`
	CurrentPeriodEnd       *time.Time         `gorm:"type:datetime"`
	CancelAtPeriodEnd      bool               `gorm:"type:tinyint(1);default:0"`

	User User             `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Plan SubscriptionPlan `gorm:"foreignKey:PlanID;constraint:OnDelete:SET NULL"`
}
