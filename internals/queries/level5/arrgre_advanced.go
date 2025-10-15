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
		UserID    uint
		AvgRating float64
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
