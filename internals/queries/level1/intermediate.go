package level1

import (
	"fmt"

	"github.com/Amanuel-0/gorm-pg/internals/database/models"
	"github.com/Amanuel-0/gorm-pg/internals/util"
	"gorm.io/gorm"
)

// Create a user and link them to multiple preferred genres.
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
		PreferredGenres: make([]*models.Genre, len(genres)),
	}

	for i, generId := range genres {
		user.PreferredGenres[i] = &models.Genre{ID: generId}
	}

	// using Traditional API
	if err := db.Create(&user).Error; err != nil {
		fmt.Printf("error creating user with genres: %v", err)
	}

	// db.Model(&models.User{}).Association("PreferredGenres").Append(genres)

	util.PrettyPrint(user, "CreateUserWithGenres: method ")
}

// Create a book and attach multiple `book_images`.
func CreateBooksWithImages(db *gorm.DB) {
	// similar question is with same kind of implementation exists
}

// Find all books with “like_new” condition by a specific author.
func GetLikeNewBooksOfAuthor(db *gorm.DB) {
	var authorId uint = 2
	// check if that author exist
	var author models.Author
	if err := db.Where("id = ?", authorId).First(&author).Error; err != nil {
		fmt.Printf("something wrong when finding author: %v", err)
	}

	var books []models.Book
	if err := db.Model(&models.Book{}).
		Where("author_id = ?", author.ID).
		// Note: `condition` is a reserved keyword for sql, so it needs to be wrapped in a
		// backtick to be used
		Where("`condition` = ?", models.ConditionLikeNew). // "like_new"
		Find(&books).Error; err != nil {
		fmt.Printf("something wrong when finding author: %v", err)

	}

	util.PrettyPrint(books, "GetLikeNewBooksOfAuthor: method")
	var count = len(books)
	fmt.Println("count: ", count)
}

// List all books in a specific city or country.
func GetBookByCityId(db *gorm.DB) {
	var city string = "San Francisco"
	var books []models.Book
	if err := db.Model(&models.Book{}).Where("location_city = ?", city).Find(&books).Error; err != nil {
		fmt.Printf("error finding books by city: %v", err)
	}
	util.PrettyPrint(books, "GetBookByCityId: method")
}

// Get all genres associated with a given book (many-to-many).
func GetGenresOfABook(db *gorm.DB) {
	var bookId uint = 11
	var genres []models.Genre

	//
	// Option 1
	//
	// if err := db.Model(&models.Genre{}).
	// 	Select("genres.id", "genres.slug", "genres.name").
	// 	// loading the many-to-many rel
	// 	Joins("JOIN book_genres bg ON bg.genre_id = genres.id").
	// 	Joins("JOIN books b ON b.id = bg.book_id").
	// 	Where("b.id = ?", bookId).
	// 	Group("genres.id").Find(&genres).Error; err != nil {
	// 	fmt.Printf("error getting genres of a book: %v", err)
	// }

	//
	// Option 2; this is the more idiomatic way of writing it in a sense that
	// 			  it is more read able and GORMs' approach
	//
	bookWithId := models.Book{ID: bookId}
	err := db.Model(&bookWithId).Association("Genres").Find(&genres)
	if err != nil {
		fmt.Printf("error getting genres of a book: %v\n", err)
		return
	}

	util.PrettyPrint(genres, "GetGenresOfABook: method")
}
