package models

import (
	"gorm.io/gorm"
)

type NotificationType string

const (
	NotificationTypeExchangeRequestReceived  NotificationType = "exchange_request_received"
	NotificationTypeExchangeRequestAccepted  NotificationType = "exchange_request_accepted"
	NotificationTypeExchangeRequestDeclined  NotificationType = "exchange_request_declined"
	NotificationTypeExchangeShipped          NotificationType = "exchange_shipped"
	NotificationTypeExchangeDelivered        NotificationType = "exchange_delivered"
	NotificationTypeExchangeCompleted        NotificationType = "exchange_completed"
	NotificationTypeNewMessageInExchange     NotificationType = "new_message_in_exchange"
	NotificationTypeNewCommunityMessage      NotificationType = "new_community_message"
	NotificationTypeBookReviewReceived       NotificationType = "book_review_received"
	NotificationTypeSubscriptionExpiringSoon NotificationType = "subscription_expiring_soon"
	NotificationTypeSubscriptionRenewed      NotificationType = "subscription_renewed"
	NotificationTypeGeneralAnnouncement      NotificationType = "general_announcement"
	// Add more as needed
)

func (nt NotificationType) IsValid() bool {
	switch nt {
	case NotificationTypeExchangeRequestReceived,
		NotificationTypeExchangeRequestAccepted,
		NotificationTypeExchangeRequestDeclined,
		NotificationTypeExchangeShipped,
		NotificationTypeExchangeDelivered,
		NotificationTypeExchangeCompleted,
		NotificationTypeNewMessageInExchange,
		NotificationTypeNewCommunityMessage,
		NotificationTypeBookReviewReceived,
		NotificationTypeSubscriptionExpiringSoon,
		NotificationTypeSubscriptionRenewed,
		NotificationTypeGeneralAnnouncement:
		return true
	}
	return false
}

type Notification struct {
	gorm.Model
	UserID  uint             `json:"user_id" gorm:"index;not null"`
	Type    NotificationType `json:"type" gorm:"size:100;not null"`
	Payload string           `json:"payload" gorm:"type:json"`
	Read    bool             `json:"read" gorm:"default:false"`

	// previous database design field -- reverse if needed
	// Type    string `json:"type" gorm:"size:100;not null"`

	// Relationships
	User *User `json:"user" gorm:"foreignKey:UserID"`
}
