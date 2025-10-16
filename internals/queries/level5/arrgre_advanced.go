package level5

import (
	"fmt"

	"github.com/Amanuel-0/gorm-pg/internals/database/models"
	"github.com/Amanuel-0/gorm-pg/internals/util"
	"gorm.io/gorm"
)

// ## ⚖️ **LEVEL 5: Aggregations, Analytics & Advanced Queries**

// Find top 5 users by total number of books owned.
func GetTop5UsersByBooksOwned(db *gorm.DB) {
	type Result struct {
		models.User
		Total int `json:"total"`
	}
	var users []Result
	if err := db.Model(&models.User{}).
		// Note: Preload here doesn't have any effect
		// Preload("UserProfile", func(db *gorm.DB) *gorm.DB {
		// 	return db.Select("id, user_id, bio")
		// }).
		Joins("JOIN books b ON b.owner_id = users.id").
		Group("users.id").
		Select("users.*, COUNT(b.owner_id) AS total").
		Order("total DESC").
		Limit(5).
		Scan(&users).Error; err != nil {
		fmt.Printf("error fetching top 5 book owner users: %v", err)
	}

	util.PrettyPrint(users, "GetTop5UsersByBooksOwned: method")
	fmt.Println("\nuser count: ", len(users))
}

// Find authors with the most books listed.
func AuthorsWithMostBookListed(db *gorm.DB) {
	type Result struct {
		models.Author
		TotalBooks uint `json:"total_books"`
	}
	var authors []Result
	res := db.Model(&models.Author{}).
		Joins("JOIN books b ON b.author_id = authors.id").
		Where("b.active = ?", true).
		Select("authors.id, authors.name, COUNT(b.id) AS total_books").
		Group("authors.id").
		Order("total_books DESC").
		Limit(10).
		Scan(&authors)

	if err := res.Error; err != nil {
		fmt.Printf("error getting authors with most listed books: %v", err)
	}

	util.PrettyPrint(authors, "AuthorsWithMostBookListed: method")
	fmt.Println("author counts: ", res.RowsAffected)
}

// Calculate the average rating per user from `user_ratings`.
func GetAvgUserRating(db *gorm.DB) {
	type Result struct {
		UserID    uint    `json:"user_id"`
		AvgRating float64 `json:"avg_rating"`
	}
	var result []Result
	res := db.Model(&models.UserRating{}).
		Select("user_ratings.rated_user_id AS user_id, AVG(user_ratings.rating) AS avg_rating").
		Group("user_ratings.rated_user_id").
		Order("avg_rating DESC").
		Scan(&result)

	util.PrettyPrint(result, "GetAvgUserRating: method")
	fmt.Println("count: ", res.RowsAffected)
}

// List books with more than or equal to 2 reviews and an average rating > 4.
func GetBooksWithCondReviewAndRaring(db *gorm.DB) {
	type Result struct {
		models.Book
		TotalReviews float64 `json:"total_review"`
		AvgRating    float64 `json:"avg_rating"`
	}

	// var books []models.Book
	var books []Result
	result := db.Model(&models.Book{}).
		Joins("JOIN book_reviews br ON br.book_id = books.id").
		Select("books.*, COUNT(br.id) AS total_reviews, AVG(br.rating) AS avg_rating").
		Group("books.id").
		Having("total_reviews >= ? AND avg_rating > ?", 2, 4).
		Order("books.id DESC").
		Preload("BookReviews", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "book_id", "reviewer_id", "rating", "comment")
		}).
		Find(&books)

	util.PrettyPrint(books, "GetBooksWithCondReviewAndRaring: method")
	fmt.Println("\ncount: ", result.RowsAffected)
}

// Count the number of messages sent per user per month.
func ListUsersMessageStatsByMonth(db *gorm.DB) {
	type Result struct {
		models.User
		// models.User `gorm:"embedded"` // to make sure fields are mapped correctly
		Year  uint `json:"year"`
		Month uint `json:"month"`
		Count uint `json:"count"`
	}
	var results []Result
	r := db.Model(&models.User{}).
		Select("users.*, YEAR(msgs.created_at) AS year, MONTH(msgs.created_at) AS month, COUNT(msgs.id) AS count").
		Joins("JOIN messages msgs ON msgs.sender_id = users.id").
		Group("users.id, year, month").
		Preload("UserProfile", func(db *gorm.DB) *gorm.DB {
			return db.Select("user_id", "id", "bio", "avatar_url")
		}).
		Order("year DESC, month DESC, users.id").
		Find(&results)

		// Note: Preload - does not work with Scan
		// Scan(&results)

	util.PrettyPrint(results, "ListUsersMessageStatsByMonth: method")
	fmt.Println("count: ", r.RowsAffected)
}

// Find users who have never participated in an exchange.
func GetUsersWithNoExchangeHistory(db *gorm.DB) {
	var users []models.User
	r := db.Model(&models.User{}).
		Joins("LEFT JOIN exchanges ex ON (ex.requester_id = users.id OR ex.responder_id = users.id)").
		Group("users.id").
		Preload("UserProfile", func(db *gorm.DB) *gorm.DB {
			return db.Select("user_id", "id", "bio")
		}).
		Where("ex.id IS NULL").
		Find(&users)

	if err := r.Error; err != nil {
		fmt.Printf("error fetching users with no exchange history: %v", err)
	}

	util.PrettyPrint(users, "GetUsersWithNoExchangeHistory: method")
	fmt.Println("count: ", r.RowsAffected)
}

// Calculate total revenue per month from `payments`.
func TotalRevenuePerMonth(db *gorm.DB) {
	type Result struct {
		Year   uint `json:"year"`
		Month  uint `json:"month"`
		Amount uint `json:"amount"` // amount is in cents
	}
	var results []Result
	r := db.Model(&models.Payment{}).
		Select("payments.*, YEAR(created_at) AS year, MONTH(created_at) AS month, SUM(amount_cents) AS amount").
		Group("year, month").
		Scan(&results)

	if r.Error != nil {
		fmt.Printf("error fetching total revenue per month: %v", r.Error)
	}

	util.PrettyPrint(results, "TotalRevenuePerMonth: method")
	fmt.Println("\ncount: ", r.RowsAffected)
}
