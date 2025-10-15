package level1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"

	// "context"

	"github.com/Amanuel-0/gorm-pg/internals/database/models"
	"github.com/Amanuel-0/gorm-pg/internals/util"
	"gorm.io/gorm"
)

// ### ✅ Simple tasks

// Create a new user with a profile.
func CreateUser(db *gorm.DB) {
	user := models.User{
		Email: "jegna@gmail.com",
		Phone: "251963158999",
		UserProfile: models.UserProfile{
			FirstName: "Amanuel",
			LastName:  "Girma",
			Bio:       "I'm Amanuel Girma. I am a Software Developer with 4+ years of experience.",
		},
	}

	// using Traditional API
	if err := db.Create(&user).Error; err != nil {
		fmt.Printf("error creating user with Traditional API: %v", err)
	}

	// using GORM Generic method
	// ctx := context.Background()
	// if err := gorm.G[models.User](db).Create(ctx, &user).Error(); err != nil {
	// 	fmt.Printf("Error creating user with Generic API: %v\n", err)
	// }
}

// Fetch a user by email.
func GetUserByEmail(db *gorm.DB, email string) {
	// using traditional API
	var user models.User
	// 1. using inline condition
	// if err := db.First(&user, "email = ?", email).Error; err != nil {
	// 	fmt.Printf("Error getting user with Traditional API: %v\n", err)
	// }
	// 2. using where chain method
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		fmt.Printf("error getting user with Traditional API: %v", err)
	}

	// print the user
	util.PrettyPrint(user, "get user by email")

	// using Generic API -- similar
}

// Update a user’s display name and avatar.
func UpdateUser(db *gorm.DB, id uint) {
	// using traditional API
	updates := map[string]interface{}{
		"FirstName": "Amanuel Updated",
		"AvatarURL": "https://example.com/new-avatar.jpg",
	}
	if err := db.Model(&models.UserProfile{}).Where("user_id = ?", id).Updates(updates).Error; err != nil {
		fmt.Printf("error updating user with Traditional API: %v", err)
	}
}

// Soft delete a user (set `deleted_at`).
func DeleteUser(db *gorm.DB, id uint) {
	// db.Where("user_id").Delete(&)
	ctx := context.Background()
	affected, err := gorm.G[models.User](db).Where("id = ?", id).Delete(ctx)
	fmt.Printf("affected: %d", affected)
	if err != nil {
		fmt.Printf("error deleting user: %v", err)
	}
}

// Create an author and a few genres.
func CreateAuthor(db *gorm.DB) {
	author := models.Author{
		// the writer of the `Game of Thrones Series`
		Name: "George RR Martin",
	}
	// create if not found, so that the query could not fail
	if err := db.Model(&models.Author{}).Where("name = ?", author.Name).FirstOrCreate(&author).Error; err != nil {
		fmt.Printf("Upsert Author: %v", author)
	}

	util.PrettyPrint(author, "Created Author")
}

// Create a book with an author and assign multiple genres.
func CreateBook(db *gorm.DB) {
	// the userID of the creator/owner of the book
	const userID uint = 1

	var author models.Author
	if err := db.
		Where("name = ?", "George RR Martin").
		First(&author).Error; err != nil {
		fmt.Printf("error getting author with Traditional API: %v", err)
	}
	util.PrettyPrint(author, "Fetched Author")

	// get 2 genres IDs
	var genres []uint
	db.Model(&models.Genre{}).Where("name IN ?", []string{"Fantasy", "Thriller"}).Pluck("ID", &genres)
	util.PrettyPrint(genres, "Plucked Genres")

	// prepare teh book object
	book := models.Book{
		Title:       "A Game of Thrones",
		Description: "A Game of Thrones is the first book in A Song of Ice and Fire, a series of fantasy novels by American author George R. R. Martin.",
		Condition:   models.ConditionLikeNew,
		OwnerID:     userID,
		AuthorID:    &author.ID,
		Active:      false,
		Genres:      make([]*models.Genre, len(genres)),
	}

	for i, genreID := range genres {
		// book.Genres[i] = models.Genre{Model: gorm.Model{ID: genreID}}
		book.Genres[i] = &models.Genre{ID: genreID}
	}

	// create the book
	if err := db.Create(&book).Error; err != nil {
		fmt.Printf("error creating books. %v", err)
	}

	util.PrettyPrint(&book, "created book data")

}

// **************************************************//
//
// START  Retrieve all books owned by a specific user.
//
// **************************************************//
func GetBooksOfUserAwesome(db *gorm.DB) {
	const userId uint = 1
	var books []models.Book
	db.Model(&models.Book{}).
		Select("id", "title", "active", "created_at", "updated_at", "owner_id", "author_id").
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Preload("Owner", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "email")
		}).
		Preload("Owner.UserProfile", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "user_id", "display_name", "bio")
		}).
		Preload("Genres", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "slug")
		}).
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "is_primary")
		}).
		Where("owner_id = ?", userId).
		Find(&books)

	// fmt.Println("Awesome List of books: ", books)
	util.PrettyPrint(books, "Awesome List of Books")

}

func GetBooksOfUser(db *gorm.DB) {
	const userID uint = 1

	// type Result struct {
	// 	Title string
	// 	Email string
	// }

	// var books []Result
	// db.Table("books").
	// 	Joins("JOIN users ON users.id = books.owner_id").
	// 	// Joins("JOIN user_profiles ON user_profiles.user_id = users.id").
	// 	Select("books.title, users.email").
	// 	Where("owner_id = ?", userID).
	// 	Scan(&books)

	// var res []map[string]interface{}
	// db.Model(&models.Book{}).
	// 	Preload("Owner").
	// 	Preload("Owner.UserProfile").
	// 	Preload("Author").
	// 	Where("owner_id = ?", userID).
	// 	Select("id", "title", "owner_id", "author_id", "condition", "created_at", "updated_at").
	// 	Scan(&res)
	// Find(&res)
	// Statement.Scan(&res)
	// Find(&res)

	// var books []Result
	// db.Model(&models.Book{}).
	// 	Preload("Owner", func(db *gorm.DB) *gorm.DB {
	// 		return db.Select("id", "email")
	// 	}).
	// 	Preload("Owner.UserProfile", func(db *gorm.DB) *gorm.DB {
	// 		return db.Select("id", "user_id", "avatar_url")
	// 	}).
	// 	Preload("Author", func(db *gorm.DB) *gorm.DB {
	// 		return db.Select("id", "name")
	// 	}).
	// 	Preload("Genres", func(db *gorm.DB) *gorm.DB {
	// 		return db.Select("id", "name")
	// 	}).
	// 	Preload(clause.Associations).
	// 	Select("id", "title", "owner_id", "author_id", "condition", "created_at", "updated_at").
	// 	Where("owner_id = ?", userID).
	// 	Scan(&books)

	type BookWithOwner struct {
		ID          uint
		Title       string
		Email       string
		OwnerAvatar string
		AuthorName  string
		// Owner struct {
		// }
		Genres []struct {
			ID   uint
			Slug string
			Name string
		}
	}

	var books []BookWithOwner
	rows, err := db.Table("books").
		Preload("Genres").
		Joins("JOIN users ON users.id = books.owner_id").
		Joins("JOIN user_profiles ON user_profiles.user_id = users.id").
		Joins("JOIN authors ON authors.id = books.author_id").
		// many-to-many relation with genres
		Joins("JOIN book_genres ON book_genres.book_id = books.id").
		Joins("JOIN genres ON genres.id = book_genres.genre_id").
		Select(`
        books.id,
        books.title,
        users.email AS email,
		user_profiles.avatar_url AS owner_avatar,
        authors.name AS author_name,
		genres.id AS id,
		genres.name AS name,
		genres.slug AS slug
    `).
		Where("books.owner_id = ?", userID).
		// Group("books.id").
		Rows()
	// Scan(&books)
	if err != nil {
		fmt.Printf("error fetching books: %v", err)
	}
	defer rows.Close()

	// cols, _ := rows.Columns()
	// fmt.Println("\ncolumns are: ", cols, "\n")

	for rows.Next() {
		var book BookWithOwner
		err := rows.Scan(util.CollectFieldPtrs(reflect.ValueOf(&book).Elem())...)
		if err != nil {
			fmt.Printf("error scanning book row: %v", err)
			continue
		}
		books = append(books, book)
	}

	util.PrettyPrint(books, "Books of User")

}

func GetBooksOfUser2(db *gorm.DB) {
	const userID uint = 1

	type GenreSummary struct {
		ID   uint
		Name string
		Slug string
	}

	type BookWithOwner struct {
		ID          uint
		Title       string
		Email       string
		OwnerAvatar string
		AuthorName  string
		Genres      []GenreSummary
	}

	var books []BookWithOwner
	booksMap := make(map[uint]*BookWithOwner)

	rows, err := db.Table("books").
		Joins("JOIN users ON users.id = books.owner_id").
		Joins("JOIN user_profiles ON user_profiles.user_id = users.id").
		Joins("JOIN authors ON authors.id = books.author_id").
		Joins("JOIN book_genres ON book_genres.book_id = books.id").
		Joins("JOIN genres ON genres.id = book_genres.genre_id").
		Select(`
        books.id,
        books.title,
        users.email AS email,
        user_profiles.avatar_url AS owner_avatar,
        authors.name AS author_name,
        genres.id AS genre_id,
        genres.name AS genre_name,
        genres.slug AS genre_slug
    `).
		Where("books.owner_id = ?", userID).
		Rows()
	if err != nil {
		log.Println("query error:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			bookID      uint
			title       string
			email       string
			ownerAvatar string
			authorName  string
			genreID     uint
			genreName   string
			genreSlug   string
		)
		if err := rows.Scan(
			&bookID, &title, &email, &ownerAvatar, &authorName,
			&genreID, &genreName, &genreSlug,
		); err != nil {
			log.Println("scan error:", err)
			continue
		}
		fmt.Println("Scanned row:", bookID, title, genreName)

		book, exists := booksMap[bookID]
		if !exists {
			book = &BookWithOwner{
				ID:          bookID,
				Title:       title,
				Email:       email,
				OwnerAvatar: ownerAvatar,
				AuthorName:  authorName,
			}
			booksMap[bookID] = book
		}

		book.Genres = append(book.Genres, GenreSummary{
			ID:   genreID,
			Name: genreName,
			Slug: genreSlug,
		})
	}

	// Convert map → slice
	for _, b := range booksMap {
		books = append(books, *b)
	}

	// util.PrettyPrint(books, "books")

}

// a function get books and the name of genres concatenated  by a comma
func GetBooksOfUserWithGenreStr(db *gorm.DB) {

	const userId uint = 1

	var books []models.Book

	db.Model(&models.Book{}).
		Preload("Genres", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "slug")
		}).
		Where("books.owner_id = ?", userId).
		Find(&books)

	// for _, book := range books {
	// util.PrettyPrint(book.Genres, "Genres")
	// fmt.Println("genre: ", )
	// }

	util.PrettyPrint(books, "List of books")

	/**
	* original code
	**/
	// type BookRes struct {
	// 	ID         uint
	// 	Title      string
	// 	AuthorName string
	// 	// GenreSlug  string
	// 	// GenreName  string
	// 	GenreNames string
	// }
	// var result []BookRes

	// db.Table("books").
	// 	Joins("JOIN users ON users.id = books.owner_id").
	// 	Joins("JOIN authors ON authors.id = books.author_id").
	// 	Joins("JOIN book_genres ON book_genres.book_id = books.id").
	// 	Joins("JOIN genres ON genres.id = book_genres.genre_id").
	// 	Select(`
	// 	books.id,
	// 	books.title,
	// 	authors.name as author_name,
	// 	GROUP_CONCAT(DISTINCT genres.name ORDER BY genres.name SEPARATOR ', ') AS genre_names
	// 	`).
	// 	// genres.name as genre_name,
	// 	// genres.slug as genre_slug
	// 	Group("books.id").
	// 	Where("books.owner_id = ?", userId).
	// 	Scan(&result)

	// util.PrettyPrint(result, "Books With Genre Str:")

}

// **************************************************//
//
// END  Retrieve all books owned by a specific user.
//
// **************************************************//

// Fetch a book including its author and genres.
func GetBookById(db *gorm.DB) {
	const bookId uint = 1
	var book models.Book
	if err := db.Model(&models.Book{}).
		Select(` id, title, author_id, owner_id `).
		// Select("id", "title", "author_id", "owner_id").
		Preload("Owner", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "email")
		}).
		Preload("Owner.UserProfile", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "user_id", "bio")
		}).
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Preload("Genres", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "slug")
		}).
		Preload("Owner.PreferredGenres", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "user_id")
		}).
		Where("id = ?", bookId).
		First(&book).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Printf("error fetching book by id: %v", err)
		}
	}

	util.PrettyPrint(book, "GetBookById: method")
}

//

// List all users who prefer a certain genre.
func GetAllUsersByPreferredGenre(db *gorm.DB) {
	// a slice of genre or a single genre
	var fltGenres []uint = []uint{1, 2, 3, 4}
	var users []models.User

	if err := db.Model(&models.User{}).
		Select("id", "email", "phone", "first_name", "last_name").
		Joins("JOIN user_preferred_genres pg ON pg.user_id = id").
		Group("id").
		Where("pg.genre_id IN ?", fltGenres).
		Preload("PreferredGenres", func(db *gorm.DB) *gorm.DB {
			return db.
				Select("id", "name", "slug")
		}).
		Preload("UserProfile", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "display_name", "bio", "user_id")
		}).
		Find(&users).Error; err != nil {
		fmt.Printf("error fetching users by there genres: %v", err)
	}

	util.PrettyPrint(users, "GetAllUsersByPreferredGenre/s")

}

// Add a new subscription plan.
func CreateSubsPlan(db *gorm.DB) {
	// Define example feature list as JSON
	features, _ := json.Marshal([]string{
		"Unlimited Projects",
		"Priority Support",
		"Custom Branding",
		"Team Collaboration Tools",
		"Advanced Analytics Dashboard",
	})

	var subp = models.SubscriptionPlan{
		Slug:        "pro-annual",
		Name:        "Pro Annual Plan",
		Description: "Best for teams and professionals who need advanced features and annual savings.",
		PriceCents:  9900, // $99.00
		Currency:    "USD",
		Interval:    "year",
		Features:    features,
		Active:      true,
	}

	// fetches if it already exists or create it
	if err := db.FirstOrCreate(&subp).Error; err != nil {
		fmt.Printf("error creating subscription plan: %v", err)
	}

	util.PrettyPrint(subp, "CreateSubsPlan: method")
}
