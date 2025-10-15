package level3

import (
	"errors"
	"fmt"
	"time"

	"github.com/Amanuel-0/gorm-pg/internals/database/models"
	"github.com/Amanuel-0/gorm-pg/internals/util"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ## 💳 **LEVEL 3: Transactions & Lifecycle**

//  - [ ] When a new user subscribes:
//    - [ ] Create a new `subscription` record.
//    - [ ] Create an initial `payment` record. --  (I don't think we have this)
//    - [ ] Update user’s role to `premium` (or store in metadata). -- (we don't have such kind of relation / but, an alternative is a subscription plan )

func ptrTime(t time.Time) *time.Time {
	return &t
}

func CreateSubscription(db *gorm.DB) {
	const userId = 1
	const subPlanId = 2

	// note: a user can only have a single active subscription

	// get the subscription plan to use it to calculate the sub end time
	var subPlan models.SubscriptionPlan
	if err := db.Model(&models.SubscriptionPlan{}).Where("id = ?", subPlanId).First(&subPlan).Error; err != nil {
		fmt.Printf("error getting the selected subscription plan: %v ", err)
	}
	// sub plan end date
	endDate := getSubscriptionEndDate(subPlan)
	db.Transaction(func(db *gorm.DB) error {
		// Serialize per-user subscription changes to avoid races
		var user models.User
		if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", userId).First(&user).Error; err != nil {
			return err
		}

		// Enforce single active subscription per user (lock matching rows)
		var existing models.Subscription
		err := db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ? AND status = ?", userId, models.SubscriptionStatusActive).First(&existing).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if existing.ID != 0 {
			return fmt.Errorf("user %d already has an active subscription (id=%d)", userId, existing.ID)
		}
		// 1. create a subscription
		var sub = models.Subscription{
			UserID:             userId,
			PlanID:             subPlanId,
			Status:             models.SubscriptionStatusActive,
			CurrentPeriodStart: ptrTime(time.Now()),
			CurrentPeriodEnd:   ptrTime(endDate),
		}
		if err := db.Create(&sub).Error; err != nil {
			return err
		}

		// 2. create a payment
		var pmt = models.Payment{
			UserID:         userId,
			SubscriptionID: sub.ID,
			AmountCents:    float64(subPlan.PriceCents),
			Status:         models.PaymentStatusSucceeded,
			Metadata:       datatypes.JSON([]byte(`{"order_id": "12345", "payment_method": "stripe"}`)),
		}
		if err := db.Create(&pmt).Error; err != nil {
			return err
		}

		return nil
	})

}

// - [ ] On book deletion:
//   - [ ] Soft delete the book (`archived_at`).
//   - [ ] Cascade delete its images and related `book_genres`.
func SoftDelBook(db *gorm.DB) {
	const id uint = 3

	db.Transaction(func(tx *gorm.DB) error {
		var book models.Book
		if err := tx.
			Preload("Genres").
			Preload("Images").
			Where("id = ?", id).
			First(&book).Error; err != nil {
			return err
		}

		// Clear many-to-many relationship manually
		if err := tx.Model(&book).Association("Genres").Clear(); err != nil {
			return err
		}

		// Delete Book (this will also delete Images because of OnDelete:CASCADE)
		if err := tx.Unscoped().Delete(&book).Error; err != nil {
			return err
		}

		return nil
	})
}

// - [ ] Create a transaction that handles an exchange:
//   - [ ] Mark exchange as `completed`.
//   - [ ] Update both books as unavailable.
//   - [ ] Insert two user ratings.
func CompleteExchange(db *gorm.DB) {
	const id uint = 3 // previously in 'accepted' state

	db.Transaction(func(db *gorm.DB) error {
		var ex models.Exchange

		if err := db.Model(&models.Exchange{}).
			Preload("RequesterBook").
			Preload("ResponderBook").
			First(&ex, "id = ?", id).Error; err != nil {
			return err
		}

		ex.Status = string(models.ExchangeStatusCompleted)

		db.Save(&ex)

		// make the books unavailable date set to nil & update the status
		var rqBookId = ex.RequesterBookID
		var rsBookId = ex.ResponderBookID
		// Update both books as unavailable
		if err := db.Model(&models.Book{}).
			Where("id IN ?", []uint{*rqBookId, *rsBookId}).
			Updates(models.Book{
				AvailableFrom:  nil,
				AvailableUntil: nil,
				Active:         true,
			}).Error; err != nil {
			return err
		}

		// create 2 user rating
		// requester rating
		var rqRating = models.UserRating{
			ExchangeID:  ex.ID,
			RaterID:     ex.RequesterID,
			RatedUserID: *ex.ResponderID,
			Rating:      4,
			Comment:     "I had a great experience with this person. The book was great reading, and it was in a great condition.",
		}
		// responder rating
		var rsRating = models.UserRating{
			ExchangeID:  ex.ID,
			RaterID:     *ex.ResponderID,
			RatedUserID: ex.RequesterID,
			Rating:      5,
			Comment:     "I had a great experience with this person. The book was great reading, and it was in a great condition.",
		}

		// create user rating in batch
		var ratings = []models.UserRating{rqRating, rsRating}
		if err := db.Create(&ratings).Error; err != nil {
			return err
		}

		util.PrettyPrint(ex, "CompleteExchange: method")

		return nil
	})

}

// - [ ] Create a function to cancel a subscription:
//   - [ ] Set `status = 'canceled'`.
//   - [ ] Update `current_period_end` and `cancel_at_period_end`.
func CancelSubscription(db *gorm.DB) {
	var subId uint = 2
	var sub = models.Subscription{ID: subId}
	if err := db.First(&sub).Error; err != nil {
		fmt.Printf("couldn't found subscription with id: %v", sub.ID)
	}
	// update the sub
	now := time.Now()
	if err := db.Model(&sub).Updates(models.Subscription{
		Status:            models.SubscriptionStatusCanceled,
		CurrentPeriodEnd:  &now,
		CancelAtPeriodEnd: true,
	}).Error; err != nil {
		fmt.Printf("error updating subscription: %v", err)
	}
}

// - [ ] Create a function that reports a user:
//   - [ ] Insert into `reports`.
//   - [ ] Create a `notification` for the admin.
func ReportUser(db *gorm.DB) {
	db.Transaction(func(tx *gorm.DB) error {

		var adminUser models.User
		if err := tx.Where("role = ?", "admin").First(&adminUser).Error; err != nil {
			return err
		}

		var report = models.Report{
			ReporterID: 4,
			TargetType: "user",
			TargetID:   5,
			HandledBy:  adminUser.ID,
			Reason:     "Inappropriate behavior",
			Metadata:   `{"details": "User sent offensive messages."}`,
		}

		if err := tx.Create(&report).Error; err != nil {
			return err
		}

		var notification = models.Notification{
			UserID:  adminUser.ID,
			Type:    models.NotificationTypeGeneralAnnouncement,
			Payload: `{"message": "A new user report has been submitted."}`,
			Read:    false,
		}
		if err := tx.Create(&notification).Error; err != nil {
			return err
		}

		util.PrettyPrint(report, "ReportUser: method")
		util.PrettyPrint(adminUser, "ReportUser: admin user")
		util.PrettyPrint(notification, "ReportUser: notification")

		return nil
	})
}

/*
//
// helper functions
//
*/
func getSubscriptionEndDate(subPlan models.SubscriptionPlan) time.Time {
	var timeInDate time.Time
	switch subPlan.Interval {
	case models.IntervalMonth:
		timeInDate = time.Now().AddDate(0, 1, 0)
	case models.Interval3Month:
		timeInDate = time.Now().AddDate(0, 3, 0)
	case models.IntervalYear:
		timeInDate = time.Now().AddDate(1, 0, 0)
	default:
		fmt.Println("a not existing sub plan time selected!!")
	}
	return timeInDate
}
