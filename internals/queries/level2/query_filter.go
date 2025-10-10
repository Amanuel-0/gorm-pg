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
