package models

import (
	"time"

	"gorm.io/gorm"
)

type Status string

const (
	StatusRequested Status = "requested"
	StatusAccepted  Status = "accepted"
	StatusDeclined  Status = "declined"
	StatusShipped   Status = "shipped"
	StatusInTransit Status = "in_transit"
	StatusDelivered Status = "delivered"
	StatusCompleted Status = "completed"
	StatusCanceled  Status = "canceled"
	StatusDisputed  Status = "disputed"
	StatusArchived  Status = "archived"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusRequested, StatusAccepted, StatusDeclined, StatusShipped, StatusInTransit, StatusDelivered, StatusCompleted, StatusCanceled, StatusDisputed, StatusArchived:
		return true
	}
	return false
}

type Exchange struct {
	gorm.Model
	RequesterID            uint      `json:"requester_id" gorm:"not null"`
	ResponderID            *uint     `json:"responder_id"`
	RequesterBookID        *uint     `json:"requester_book_id"`
	ResponderBookID        *uint     `json:"responder_book_id"`
	Status                 string    `json:"status" gorm:"type:enum('requested','accepted','declined','shipped','in_transit','delivered','completed','canceled','disputed','archived');not null;default:'requested'"`
	RequestedAt            time.Time `json:"requested_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	StatusUpdatedAt        time.Time `json:"status_updated_at" gorm:"not null;default:CURRENT_TIMESTAMP;autoUpdateTime"`
	AgreedStartDate        time.Time `json:"agreed_start_date"`
	AgreedEndDate          time.Time `json:"agreed_end_date"`
	ShippingRequired       bool      `json:"shipping_required" gorm:"default:true"`
	ShippingProvider       string    `json:"shipping_provider"`
	ShippingTrackingNumber string    `json:"shipping_tracking_number"`
	ShippingCostCents      int       `json:"shipping_cost_cents"`
	ShippingPayerUserID    uint      `json:"shipping_payer_user_id"`
	CompletedAt            time.Time `json:"completed_at"`
	CanceledAt             time.Time `json:"canceled_at"`
	DisputeReason          string    `json:"dispute_reason" gorm:"type:text"`
	DisputeOpenedAt        time.Time `json:"dispute_opened_at"`
	Archived               bool      `json:"archived" gorm:"default:false"`
	Metadata               string    `json:"metadata" gorm:"type:json"`

	// Relationships
	Requester     User         `json:"requester,omitempty" gorm:"foreignKey:RequesterID;constraint:OnDelete:CASCADE"`
	Responder     User         `json:"responder,omitempty" gorm:"foreignKey:ResponderID;constraint:OnDelete:SET NULL"`
	RequesterBook Book         `json:"requester_book,omitempty" gorm:"foreignKey:RequesterBookID;constraint:OnDelete:SET NULL"`
	ResponderBook Book         `json:"responder_book,omitempty" gorm:"foreignKey:ResponderBookID;constraint:OnDelete:SET NULL"`
	ShippingPayer User         `json:"shipping_payer,omitempty" gorm:"foreignKey:ShippingPayerUserID;constraint:OnDelete:SET NULL"`
	ChatThreads   []ChatThread `json:"chat_threads,omitempty" gorm:"foreignKey:ExchangeID;constraint:OnDelete:CASCADE"`
	UserRatings   []UserRating `json:"user_ratings,omitempty" gorm:"foreignKey:ExchangeID;constraint:OnDelete:SET NULL"`
}
