package level2

import (
	"fmt"
	"time"

	"github.com/Amanuel-0/gorm-pg/internals/database/models"
	"github.com/Amanuel-0/gorm-pg/internals/util"
	"gorm.io/gorm"
)

// ## 🧭 **LEVEL 2: Querying & Filtering**

// Fetch all books available between two dates.
func GetBookBetweenDates(db *gorm.DB) {
	var startDate time.Time = time.Date(2025, time.October, 10, 0, 0, 0, 0, time.UTC)
	var endDate time.Time = time.Date(2025, time.October, 12, 0, 0, 0, 0, time.UTC)

	fmt.Printf("\nstart date: %v\n", startDate)
	fmt.Printf("\nend date: %v\n\n", endDate)

	var books []models.Book
	if err := db.Model(&models.Book{}).
		Where("available_from <= ?", startDate).
		Where("available_until >= ?", endDate).
		Preload("Owner", func(db *gorm.DB) *gorm.DB {
			return db.Select("id")
		}).
		Preload("Owner.UserProfile", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "user_id", "bio")
		}).
		Select("id", "title", "owner_id", "available_from", "available_until").
		Find(&books).Error; err != nil {
	}

	util.PrettyPrint(books, "GetBookBetweenDates: method")

	var count = len(books)
	fmt.Println("count: ", count)
}

// Find all users who haven’t verified their email.
// passed - no enough information, and it seems it is repetitive
func GetUsersWithVerifiedEmail(db *gorm.DB) {}

// Retrieve all active subscriptions and their plans.
func GetActiveSubsWithPlan(db *gorm.DB) {
	var subs []models.Subscription
	if err := db.Model(&models.Subscription{}).
		Where("status = ?", models.SubscriptionStatusActive).
		Preload("Plan", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "price_cents")
		}).
		Select("id", "user_id", "plan_id", "status").
		Find(&subs).Error; err != nil {

		fmt.Printf("error fetching active subs: %v", err)
	}

	util.PrettyPrint(subs, "GetActiveSubsWithPlan")
}

// Find all users whose subscription is expired.
func GetUsersWithExpiredSub(db *gorm.DB) {
	var users []models.User

	//
	// Option 1
	//
	subQuery := db.Table("subscriptions").
		Select("user_id").
		Where("current_period_end < ?", time.Now())

	if err := db.Model(&models.User{}).
		Preload("UserProfile", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "user_id", "bio")
		}).
		Where("id IN (?)", subQuery).
		Select("id", "email", "phone", "first_name", "last_name").
		Find(&users).Error; err != nil {
		fmt.Printf("error fetching users with expired subscriptions: %v", err)
	}

	util.PrettyPrint(users, "Expired Subscriptions:")

	//
	// Option 2
	//
	// 	var users []models.User

	// 	// Subquery: Get latest subscription per user
	// 	latestSubQuery := db.Table("subscriptions AS s1").
	// 		Select("s1.user_id, s1.current_period_end").
	// 		Joins(`
	//             LEFT JOIN subscriptions AS s2
	//             ON s1.user_id = s2.user_id
	//             AND s1.current_period_end < s2.current_period_end
	//         `).
	// 		Where("s2.user_id IS NULL") // ensures only the latest record per user

	// 	if err := db.
	// 		Table("users").
	// 		Joins("JOIN (?) AS subs ON subs.user_id = users.id", latestSubQuery).
	// 		// Expired if current_period_end < NOW() OR is NULL
	// 		Where("(subs.current_period_end IS NULL OR subs.current_period_end < ?)", time.Now()).
	// 		// Only normal user accounts
	// 		Where("users.role = ?", "user").
	// 		// Optional: Only active users
	// 		Where("users.is_active = ?", true).
	// 		// Select only identifiers (customize if needed)
	// 		Select("users.id", "users.email", "users.first_name", "users.last_name").
	// 		Find(&users).Error; err != nil {

	// 		fmt.Printf("error fetching users with expired subscriptions: %v", err)
	// 		return
	// 	}

	// util.PrettyPrint(users, "Expired Subscriptions:")
}

// List all users who have ever made a successful payment.
// passed - no enough information, and it seems it is repetitive
func GetUsersWithSuccessfulPayment(db *gorm.DB) {}

// Count how many books each user owns (optimized with a single query).
func GetUsersWithBookCount(db *gorm.DB) {
	type Result struct {
		UserID    uint   `json:"user_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		BookCount int    `json:"book_count"`
	}

	var results []Result
	// This subquery calculates the book count for a given user.
	// GORM will correlate `books.owner_id` with `users.id` automatically.
	subQuery := db.Model(&models.Book{}).Select("count(id)").Where("books.owner_id = users.id")

	// The main query starts from the User model for better type safety.
	err := db.Model(&models.User{}).
		Select("users.id as user_id, users.first_name, users.last_name, (?) as book_count", subQuery).
		Scan(&results).Error

	if err != nil {
		fmt.Printf("error fetching users book count: %v\n", err)
		return
	}

	util.PrettyPrint(results, "GetUsersBookCount: method")

}

// Retrieve all books with their review averages.
func GetBooksWithAvgReview(db *gorm.DB) {
	type Result struct {
		BookID    uint    `json:"book_id"`
		Title     string  `json:"title"`
		AvgReview float64 `json:"avg_review"`
	}
	var results []Result

	subQuery := db.Model(&models.BookReview{}).Select("AVG(rating)").Where("book_reviews.book_id = books.id")

	err := db.Model(&models.Book{}).
		Select("books.id as book_id, books.title, (?) as avg_review", subQuery).
		Scan(&results).Error

	if err != nil {
		fmt.Printf("error fetching books with avg review: %v\n", err)
		return
	}

	util.PrettyPrint(results, "GetUsersBookCount: method")
}

// List all exchanges in the “requested” state with book and user info.
func GetExchangesWithRequestedStatus(db *gorm.DB) {
	var exchanges []models.Exchange

	err := db.Model(&models.Exchange{}).
		Preload("Requester", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "email", "phone")
		}).
		Preload("Requester.UserProfile", func(db *gorm.DB) *gorm.DB {
			return db.Select("user_id", "id", "bio")
		}).
		Preload("Responder", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "email", "phone")
		}).
		Preload("Responder.UserProfile", func(db *gorm.DB) *gorm.DB {
			return db.Select("user_id", "id", "bio")
		}).
		Where("status = ?", models.ExchangeStatusRequested).
		// Select("").
		Find(&exchanges).Error

	if err != nil {
		fmt.Printf("error fetching exchanges with requested status: %v", err)
	}

	util.PrettyPrint(exchanges, "GetExchangesWithRequestedStatus: method")
}

// Fetch messages in a chat thread, sorted newest → oldest.
func GetThreadMessagesSorted(db *gorm.DB) {
	// chat thread and messages have a one-to-many relationship
	var threadId uint = 2
	var messages []models.Message
	err := db.Model(&models.Message{}).
		Where("thread_id = ?", threadId).
		Preload("Thread", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "exchange_id")
		}).
		// Preload("Thread.Exchange").
		Preload("Thread.Exchange", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "requester_id", "responder_id")
		}).
		Preload("Thread.Exchange.Requester", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "email", "first_name", "last_name")
		}).
		Preload("Thread.Exchange.Responder", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "email", "first_name", "last_name")
		}).
		Select("id", "thread_id", "sender_id", "type", "body", "attachments", "created_at", "updated_at").
		Order("created_at DESC").
		Find(&messages).Error

	if err != nil {
		fmt.Printf("error fetching thread messages: %v", err)
	}

	util.PrettyPrint(messages, "GetThreadMessagesSorted: method")
}

// Find all users who haven’t logged in for 30+ days.
func GetUsersInActiveForOverAMonth(db *gorm.DB) {
	var userIDs []uint
	db.Model(&models.ActivityLog{}).
		Where("action = ?", models.LogActionLogin).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -30)).
		Pluck("user_id", &userIDs)

	fmt.Print("user ids: ", userIDs)

	var inactiveUsers []models.User
	if err := db.Model(&models.User{}).
		Where("id NOT IN ?", userIDs).
		Preload("UserProfile", func(db *gorm.DB) *gorm.DB {
			return db.Select("user_id", "id", "bio", "display_name")
		}).
		Select("id", "email", "first_name", "last_name").
		Find(&inactiveUsers).Error; err != nil {
		fmt.Printf("error fetching inactive users: %v\n", err)
	}

	util.PrettyPrint(inactiveUsers, "GetUsersInActiveForOverAMonth: method")

}
