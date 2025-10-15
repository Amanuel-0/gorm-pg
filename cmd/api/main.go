package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/Amanuel-0/gorm-pg/internals/config"
	"github.com/Amanuel-0/gorm-pg/internals/database"
	"github.com/Amanuel-0/gorm-pg/internals/database/models"
	"github.com/Amanuel-0/gorm-pg/internals/queries/level5"

	// "github.com/Amanuel-0/gorm-pg/internals/queries/level1"

	"github.com/Amanuel-0/gorm-pg/internals/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	fmt.Println("Hello, API!")

	// Load Configuration
	config, err := config.New()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// database connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DB.Username, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.DBName)
	db, err := database.ConnectDB(dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	// Migrate database tables
	autoMigrateTables(db)
	//
	//
	// Seed initial data (idempotent)
	// if err := seeder.SeedAll(db); err != nil {
	// 	log.Fatalf("failed to seed database: %v", err)
	// }
	//
	//

	// Queries
	//
	// level 1
	//
	// simple_tasks
	//
	// level1.CreateAuthor(db)
	// level1.CreateBook(db)
	// level1.GetBooksOfUser2(db)
	// level1.GetBooksOfUserWithGenreStr(db)
	// level1.GetBooksOfUserAwesome(db)
	// level1.GetBookById(db)
	// level1.GetAllUsersByPreferredGenre(db)
	// level1.CreateSubsPlan(db)

	//
	// intermediate
	//
	// level1.CreateUserWithGenres(db)
	// level1.GetLikeNewBooksOfAuthor(db)
	// level1.GetBookByCityId(db)
	// level1.GetGenresOfABook(db)

	//
	// Level 2
	//
	// level2.GetBookBetweenDates(db)
	// level2.GetActiveSubsWithPlan(db)
	// level2.GetUsersWithExpiredSub(db)
	// level2.GetUsersWithBookCount(db)
	// level2.GetBooksWithAvgReview(db)
	// level2.GetExchangesWithRequestedStatus(db)
	// level2.GetThreadMessagesSorted(db)
	// level2.GetUsersInActiveForOverAMonth(db)

	//
	// Level 3
	//
	// level3.CreateSubscription(db)
	// level3.SoftDelBook(db)
	// level3.CompleteExchange(db)
	// level3.CancelSubscription(db)
	// level3.ReportUser(db)

	//
	// Level 4
	//
	// level4.GetChatThreadsOfExchange(db)
	// level4.GetCommunityThreads(db)
	// level4.GetUsersWithWithAtLeast2Communities(db)
	// level4.GetPaidCommunities(db)
	// level4.GetExchangesOfUser(db)
	// level4.GetBooksOfUserInCompletedExchanges(db)

	//
	// Level 5
	//
	level5.GetTop5UsersByBooksOwned(db)

	// playWithGORMqueries(db)

}

func autoMigrateTables(db *gorm.DB) {
	if err := db.AutoMigrate(
		&models.ActivityLog{},
		&models.Author{},
		&models.BookImage{},
		&models.BookReview{},
		&models.Book{},
		&models.ChatThread{},
		&models.CommunityMember{},
		&models.CommunityMessage{},
		&models.CommunityThread{},
		&models.Community{},
		&models.Exchange{},
		&models.Genre{},
		// locations starts
		&models.Country{},
		&models.State{},
		&models.City{},
		// locations ends
		&models.MessageQuotaUsage{},
		&models.Message{},
		&models.ModerationAction{},
		&models.Notification{},
		&models.Report{},
		&models.SubscriptionPlan{},
		&models.Subscription{},
		&models.UserProfile{},
		&models.UserRating{},
		&models.User{},
		&models.Payment{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}

// legacy seeder function removed in favor of seeder.SeedAll

// a function that can be used to marshal and print the output
// to a terminal in a pretty format
func prettyPrint(data interface{}, message string) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal data: %v", err)
	}
	fmt.Printf("%s: %s\n", message, string(jsonData))
}

// collectFieldPtrs recursively walks a struct and collects pointers to fields.

// Play with GORM database queries
func playWithGORMqueries(db *gorm.DB) {
	// ctx := context.Background()
	// user, err := gorm.G[models.User](db).First(ctx)
	// method 1
	var user models.User
	ur := db.First(&user)
	if ur.Error != nil {
		log.Fatalf("failed to get user: %v", ur.Error)
	}
	// prettyPrint(user, "Method 1: User")

	// method 2
	// result := map[string]interface{}{}
	// r := db.Model(&models.User{}).First(&result)
	// count := r.RowsAffected
	// fmt.Printf("Total users: %d\n", count)
	// if r.Error != nil {
	// 	log.Fatalf("failed to get user: %v", r.Error)
	// }
	// prettyPrint(result, "Method 2: User")

	// row scan a user and set user
	var userScan struct {
		ID          uint   `json:"id"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		UserProfile struct {
			Bio       string  `json:"bio"`
			AvatarURL string  `json:"avatar_url"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Linkedin  string  `json:"linkedin"`
		} `json:"user_profile"`
	}

	rs := db.Table("users").
		Select("users.id, users.first_name, users.last_name, user_profiles.bio, user_profiles.avatar_url, user_profiles.latitude, user_profiles.longitude, user_profiles.linkedin").
		Joins("left join user_profiles on user_profiles.user_id = users.id").
		Row().
		Scan(util.CollectFieldPtrs(reflect.ValueOf(&userScan).Elem())...)
	if rs != nil {
		log.Fatalf("failed to scan user: %v", rs)
	}
	prettyPrint(userScan, "Row Scan User")

	// find all users
	users := []map[string]interface{}{}
	usr := db.Clauses(clause.Locking{
		Strength: "UPDATE",
	}).Model(&models.User{}).
		// Where("email LIKE ? AND first_name LIKE ? OR last_name LIKE ?", "%com%", "%i%", "%m%").
		// Select("id, first_name, last_name, email, phone").
		// Joins("left join user_profiles on user_profiles.user_id = users.id").
		Limit(10).
		// Order("id desc").
		// Distinct("phone", "email").
		Order("id desc").
		Find(&users)
	// Scan(&users)
	total := usr.RowsAffected
	fmt.Printf("Total users: %d\n", total)
	if usr.Error != nil {
		log.Fatalf("failed to get users: %v", usr.Error)
	}
	// prettyPrint(users, "Method 3: Users")

	// find a single user with a join
	var result struct {
		ID        uint   `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Bio       string `json:"bio"`
		AvatarURL string `json:"avatar_url"`
	}
	// create a struct to create a wrapper for user and user profile table to be nested
	var resultWrapper struct {
		ID          uint   `json:"id"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		UserProfile struct {
			Bio       string `json:"bio"`
			AvatarURL string `json:"avatar_url"`
		}
	}
	// var joinUser models.User
	jur := db.Model(&models.User{}).
		Select("users.id, users.first_name, users.last_name, user_profiles.bio, user_profiles.avatar_url").
		Joins("left join user_profiles on user_profiles.user_id = users.id").
		First(&result)
	if jur.Error != nil {
		log.Fatalf("failed to get joined user: %v", jur.Error)
	}

	resultWrapper.ID = result.ID
	resultWrapper.FirstName = result.FirstName
	resultWrapper.LastName = result.LastName
	resultWrapper.UserProfile.Bio = result.Bio
	resultWrapper.UserProfile.AvatarURL = result.AvatarURL

	// prettyPrint(result, "Joined User")

}
