package level4

import (
	"fmt"

	"github.com/Amanuel-0/gorm-pg/internals/database/models"
	"github.com/Amanuel-0/gorm-pg/internals/util"
	"gorm.io/gorm"
)

// ## 💬 **LEVEL 4: Relations Across Domains**

// Get all chat threads for a given exchange, including messages and senders.
func GetChatThreadsOfExchange(db *gorm.DB) {
	fmt.Println("here we care!!!")
	var exchangeId uint = 1
	var threads models.ChatThread
	if err := db.Model(&models.ChatThread{}).
		Preload("Messages").
		Preload("Creator", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "email", "phone", "first_name", "last_name", "is_active", "role")
		}).
		Preload("Creator.UserProfile", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "user_id", "bio", "avatar_url")
		}).
		Where("exchange_id = ?", exchangeId).
		Find(&threads).Error; err != nil {
		fmt.Printf("error fetching chat threads of an exchange: %v", err)
	}

	util.PrettyPrint(threads, "GetChatThreadsOfExchange: method")

}
