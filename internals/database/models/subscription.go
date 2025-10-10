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
	// gorm.Model
	ID                     uint               `json:"id,omitempty" gorm:"primaryKey"`
	UserID                 uint               `json:"user_id,omitempty" gorm:"not null;index"`
	PlanID                 uint               `json:"plan_id,omitempty" gorm:"index"`
	ProviderSubscriptionID string             `json:"provider_subscription_id,omitempty" gorm:"size:255;index"`
	Status                 SubscriptionStatus `json:"status,omitempty" gorm:"type:enum('active','past_due','canceled','trialing','expired');not null;default:'trialing'"`
	CurrentPeriodStart     *time.Time         `json:"current_period_start,omitempty" gorm:"type:datetime"`
	CurrentPeriodEnd       *time.Time         `json:"current_period_end,omitempty" gorm:"type:datetime"`
	CancelAtPeriodEnd      bool               `json:"cancel_at_period_end,omitempty" gorm:"type:tinyint(1);default:0"`

	// timestamps
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`

	// relationship
	User *User            `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Plan SubscriptionPlan `json:"plan,omitempty" gorm:"foreignKey:PlanID;constraint:OnDelete:SET NULL"`
}
