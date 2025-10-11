package models

import (
	"time"

	"gorm.io/gorm"
)

type Status string

const (
	ExchangeStatusRequested Status = "requested"
	ExchangeStatusAccepted  Status = "accepted"
	ExchangeStatusDeclined  Status = "declined"
	ExchangeStatusShipped   Status = "shipped"
	ExchangeStatusInTransit Status = "in_transit"
	ExchangeStatusDelivered Status = "delivered"
	ExchangeStatusCompleted Status = "completed"
	ExchangeStatusCancelled Status = "canceled"
	ExchangeStatusInDispute Status = "disputed"
	ExchangeStatusArchived  Status = "archived"
)

func (s Status) IsValid() bool {
	switch s {
	case ExchangeStatusRequested, ExchangeStatusAccepted, ExchangeStatusDeclined, ExchangeStatusShipped, ExchangeStatusInTransit, ExchangeStatusDelivered, ExchangeStatusCompleted, ExchangeStatusCancelled, ExchangeStatusInDispute, ExchangeStatusArchived:
		return true
	}
	return false
}

type Exchange struct {
	// gorm.Model
	ID                     uint       `json:"id,omitempty" gorm:"primaryKey"`
	RequesterID            uint       `json:"requester_id,omitempty" gorm:"not null"`
	ResponderID            *uint      `json:"responder_id,omitempty"`
	RequesterBookID        *uint      `json:"requester_book_id,omitempty"`
	ResponderBookID        *uint      `json:"responder_book_id,omitempty"`
	Status                 string     `json:"status,omitempty" gorm:"type:enum('requested','accepted','declined','shipped','in_transit','delivered','completed','canceled','disputed','archived');not null;default:'requested'"`
	RequestedAt            *time.Time `json:"requested_at,omitempty" gorm:"not null;default:CURRENT_TIMESTAMP"`
	StatusUpdatedAt        *time.Time `json:"status_updated_at,omitempty" gorm:"not null;default:CURRENT_TIMESTAMP;autoUpdateTime"`
	AgreedStartDate        *time.Time `json:"agreed_start_date,omitempty"`
	AgreedEndDate          *time.Time `json:"agreed_end_date,omitempty"`
	ShippingRequired       bool       `json:"shipping_required,omitempty" gorm:"default:true"`
	ShippingProvider       string     `json:"shipping_provider,omitempty"`
	ShippingTrackingNumber string     `json:"shipping_tracking_number,omitempty"`
	ShippingCostCents      int        `json:"shipping_cost_cents,omitempty"`
	ShippingPayerUserID    uint       `json:"shipping_payer_user_id,omitempty"`
	CompletedAt            *time.Time `json:"completed_at,omitempty"`
	CanceledAt             *time.Time `json:"canceled_at,omitempty"`
	DisputeReason          string     `json:"dispute_reason,omitempty" gorm:"type:text"`
	DisputeOpenedAt        *time.Time `json:"dispute_opened_at,omitempty"`
	Archived               bool       `json:"archived,omitempty" gorm:"default:false"`
	Metadata               *string    `json:"metadata,omitempty" gorm:"type:json"`

	// timestamps
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`

	// Relationships
	Requester     *User         `json:"requester,omitempty" gorm:"foreignKey:RequesterID;constraint:OnDelete:CASCADE"`
	Responder     *User         `json:"responder,omitempty" gorm:"foreignKey:ResponderID;constraint:OnDelete:SET NULL"`
	RequesterBook *Book         `json:"requester_book,omitempty" gorm:"foreignKey:RequesterBookID;constraint:OnDelete:SET NULL"`
	ResponderBook *Book         `json:"responder_book,omitempty" gorm:"foreignKey:ResponderBookID;constraint:OnDelete:SET NULL"`
	ShippingPayer *User         `json:"shipping_payer,omitempty" gorm:"foreignKey:ShippingPayerUserID;constraint:OnDelete:SET NULL"`
	ChatThreads   []*ChatThread `json:"chat_threads,omitempty" gorm:"foreignKey:ExchangeID;constraint:OnDelete:CASCADE"`
	UserRatings   []*UserRating `json:"user_ratings,omitempty" gorm:"foreignKey:ExchangeID;constraint:OnDelete:SET NULL"`
}
