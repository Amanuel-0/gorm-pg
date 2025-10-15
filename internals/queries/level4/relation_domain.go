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
		First(&threads).Error; err != nil {
		fmt.Printf("error fetching chat threads of an exchange: %v", err)
	}

	util.PrettyPrint(threads, "GetChatThreadsOfExchange: method")

}

// Fetch all community threads and their messages (with author/creator info).
func GetCommunityThreads(db *gorm.DB) {
	var threads []models.CommunityThread
	if err := db.Model(&models.CommunityThread{}).
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "thread_id", "sender_id", "body")
		}).
		Preload("Messages.Sender", func(db *gorm.DB) *gorm.DB {
			return db.Select("id")
		}).
		Preload("Messages.Sender.UserProfile", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "user_id", "bio", "avatar_url")
		}).
		Preload("Creator", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "email", "phone", "first_name", "last_name", "role")
		}).
		Preload("Creator.UserProfile", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "user_id", "bio", "avatar_url")
		}).
		Select("id", "community_id", "title", "created_by").
		Find(&threads).Error; err != nil {
		fmt.Printf("error fetching community threads: %v", err)
	}

	util.PrettyPrint(threads, "GetCommunityThreads: method")
}

// List all users who belong to at least 2 communities.
func GetUsersWithWithAtLeast2Communities(db *gorm.DB) {
	// 1. my first implementation
	// var users []models.User
	// var ids []uint
	// if err := db.Model(&models.CommunityMember{}).
	// 	Select("user_id", "count(user_id) as count").
	// 	Group("user_id").
	// 	Having("count(user_id) >= ?", 4).
	// 	Pluck("user_id", &ids).Error; err != nil {
	// 	fmt.Printf("error fetching user_id s: %v", err)
	// }
	// if err := db.Model(&models.User{}).
	// 	Where("id IN (?)", ids).
	// 	Find(&users).Error; err != nil {
	// 	fmt.Printf("error fetching user: %v", err)
	// }

	// 2. better way of implementing it without having an intermediate value like
	// `ids`
	var users []models.User
	if err := db.Model(&models.User{}).
		Preload("UserProfile", func(db *gorm.DB) *gorm.DB {
			return db.Select("user_id", "id", "display_name", "bio", "avatar_url")
		}).
		Select("users.*").
		Joins("JOIN community_members cm ON cm.user_id = users.id").
		Group("users.id").
		Having("COUNT(cm.community_id) >= ?", 4).
		Find(&users).Error; err != nil {
		fmt.Printf("error fetching user: %v", err)
	}

	util.PrettyPrint(users, "GetUsersWithWithAtLeast2Communities: method")
	fmt.Println("user count: ", len(users))
}
