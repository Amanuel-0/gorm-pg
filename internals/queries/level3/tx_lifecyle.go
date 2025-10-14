package level3

import (
	"errors"
	"fmt"
	"time"

	"github.com/Amanuel-0/gorm-pg/internals/database/models"
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
