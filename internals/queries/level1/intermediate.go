package level1

import (
	"fmt"

	"github.com/Amanuel-0/gorm-pg/internals/database/models"
	"github.com/Amanuel-0/gorm-pg/internals/util"
	"gorm.io/gorm"
)

func CreateUserWithGenres(db *gorm.DB) {
	var genres []uint = []uint{1, 2, 3}
	user := models.User{
		Email: "chala@gmail.com",
		Phone: "2519631589991",
		UserProfile: models.UserProfile{
			FirstName: "Chala",
			LastName:  "Chelchesa",
			Bio:       "I'm Chala Chelchesa.",
		},
		PreferredGenres: make([]models.Genre, len(genres)),
	}

	for i, generId := range genres {
		user.PreferredGenres[i] = models.Genre{ID: generId}
	}

	// using Traditional API
	if err := db.Create(&user).Error; err != nil {
		fmt.Printf("error creating user with genres: %v", err)
	}

	// db.Model(&models.User{}).Association("PreferredGenres").Append(genres)

	util.PrettyPrint(user, "CreateUserWithGenres: method ")
}
