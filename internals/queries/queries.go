package queries

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Amanuel-0/gorm-pg/internals/database/models"
	"gorm.io/gorm"
)

// Dtos

type UserDto struct {
	Email string
	Phone string
}

func GetUser(db *gorm.DB) {
	ctx := context.Background()
	// user, err := gorm.G[models.User](db).Select("email", "phone").Find(ctx)
	// if err != nil {
	// 	log.Fatalf("failed to get user: %v", err)
	// }
	// userJSON, err := json.MarshalIndent(user, "", "  ")
	// if err != nil {
	// 	log.Fatalf("failed to marshal user: %v", err)
	// }
	// fmt.Println("First User: ", string(userJSON))

	// users
	// user, err := gorm.G[models.User](db).Select("email", "phone").First(ctx)
	// // db.WithContext(ctx).Select("email", "phone").Find(&users)
	// if err != nil {
	// 	log.Fatalf("failed to get user: %v", err)
	// }
	// userJSON, err := json.MarshalIndent(user, "", "  ")
	// if err != nil {
	// 	log.Fatalf("failed to marshal user: %v", err)
	// }
	// fmt.Println("First User: ", string(userJSON))

	//  user
	var user UserStriped
	// db.WithContext(ctx).
	// 	Model(&models.User{}).
	// 	Joins("left join user_profiles on user_profiles.user_id = users.id").
	// 	Select("email", "phone", "user_profiles.id", "user_profiles.bio").
	// 	First(&user)
	// db.WithContext(ctx).
	// 	Model(&models.User{}).
	// 	Joins("left join user_profiles on user_profiles.user_id = users.id").
	// 	Select("users.email as email", "users.phone as phone", "user_profiles.bio as user_profile__bio").
	// 	First(&user)
	// err := db.WithContext(ctx).
	// 	Model(&models.User{}).
	// 	Joins("left join user_profiles on user_profiles.user_id = users.id").
	// 	Select(`
	//     users.email as email,
	//     users.phone as phone,
	//     user_profiles.bio as user_profile__bio
	// `).
	// 	First(&user).Error

	// another idiomatic way of getting nested table data
	err := db.WithContext(ctx).
		Model(&models.User{}).
		// Select("email", "phone").
		// Preload("UserProfile", func(db *gorm.DB) *gorm.DB {
		// 	return db.Select("bio", "user_id", "avatar_url")
		// }).
		First(&user).Error
	if err != nil {
		log.Fatalf("failed to get user: %v", err)
	}
	userJSON, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal user: %v", err)
	}
	println("First User: ", string(userJSON))

	// user profile
	// var userProfile models.UserProfile
	// err = db.WithContext(ctx).Model(&models.UserProfile{}).First(&userProfile).Error
	// if err != nil {
	// 	log.Fatalf("failed to get user: %v", err)

	// }
	// userProfileJSON, err := json.MarshalIndent(userProfile, "", "  ")
	// if err != nil {
	// 	log.Fatalf("failed to marshal user: %v", err)
	// }
	// println("First User Profile: ", string(userProfileJSON))

}

// type UserStriped struct {
// 	Email       string
// 	Phone       string
// 	UserProfile UserProfileStriped `gorm:"-"`
// }

// type UserProfileStriped struct {
// 	Bio string
// }

// structs for the idiomatic way of getting nested table data
type UserStriped struct {
	Email       string             `json:"email"`
	Phone       string             `json:"phone"`
	UserProfile UserProfileStriped `json:"userProfile" gorm:"foreignKey:UserID"` // pointer so it's nil if no profile exists
}

type UserProfileStriped struct {
	UserID uint `json:"-"` // needed for GORM to map foreignKey

	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}
