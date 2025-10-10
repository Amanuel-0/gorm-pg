package seeder

import (
	"context"
	"fmt"
	"time"

	"github.com/Amanuel-0/gorm-pg/internals/database/models"
	"gorm.io/gorm"
)

// SeedAll seeds the database with comprehensive relational data for all practice queries.
// It is idempotent and creates data for all 7 levels of practice scenarios.
func SeedAll(db *gorm.DB) error {
	ctx := context.Background()

	// 1. LOCATIONS - Multiple countries, states, cities for geo-based queries
	if err := seedLocations(db, ctx); err != nil {
		return fmt.Errorf("failed to seed locations: %w", err)
	}

	// 2. USERS - Diverse users with different roles, subscription statuses, and profiles
	users, err := seedUsers(db, ctx)
	if err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	// 3. AUTHORS & GENRES - Extensive collection for book relationships
	authors, genres, err := seedAuthorsAndGenres(db, ctx)
	if err != nil {
		return fmt.Errorf("failed to seed authors and genres: %w", err)
	}

	// 4. BOOKS - Books with varied conditions, locations, and relationships
	books, err := seedBooks(db, ctx, users, authors, genres)
	if err != nil {
		return fmt.Errorf("failed to seed books: %w", err)
	}

	// 5. SUBSCRIPTIONS - Multiple plans and various subscription statuses
	_, err = seedSubscriptions(db, ctx, users)
	if err != nil {
		return fmt.Errorf("failed to seed subscriptions: %w", err)
	}

	// 6. COMMUNITIES - Multiple communities with different settings
	communities, err := seedCommunities(db, ctx, users)
	if err != nil {
		return fmt.Errorf("failed to seed communities: %w", err)
	}

	// 7. EXCHANGES - Exchanges in various states for testing
	exchanges, err := seedExchanges(db, ctx, users, books)
	if err != nil {
		return fmt.Errorf("failed to seed exchanges: %w", err)
	}

	// 8. MESSAGES & CHAT - Communication data
	if err := seedMessages(db, ctx, users, exchanges, communities); err != nil {
		return fmt.Errorf("failed to seed messages: %w", err)
	}

	// 9. REVIEWS & RATINGS - For analytics and user feedback
	if err := seedReviewsAndRatings(db, ctx, users, books, exchanges); err != nil {
		return fmt.Errorf("failed to seed reviews and ratings: %w", err)
	}

	// 10. NOTIFICATIONS - Various notification types
	if err := seedNotifications(db, ctx, users); err != nil {
		return fmt.Errorf("failed to seed notifications: %w", err)
	}

	// 11. MODERATION - Reports and moderation actions
	if err := seedModeration(db, ctx, users); err != nil {
		return fmt.Errorf("failed to seed moderation: %w", err)
	}

	// 12. ACTIVITY LOGS - For audit trails
	if err := seedActivityLogs(db, ctx, users, books); err != nil {
		return fmt.Errorf("failed to seed activity logs: %w", err)
	}

	// 13. MESSAGE QUOTA USAGE - For usage tracking
	if err := seedMessageQuotaUsage(db, ctx, users); err != nil {
		return fmt.Errorf("failed to seed message quota usage: %w", err)
	}

	return nil
}

// seedLocations creates diverse geographic data for geo-based queries
func seedLocations(db *gorm.DB, ctx context.Context) error {
	countries := []models.Country{
		{Code: "US", Name: "United States"},
		{Code: "CA", Name: "Canada"},
		{Code: "GB", Name: "United Kingdom"},
		{Code: "AU", Name: "Australia"},
		{Code: "DE", Name: "Germany"},
		{Code: "FR", Name: "France"},
		{Code: "JP", Name: "Japan"},
		{Code: "BR", Name: "Brazil"},
		{Code: "IN", Name: "India"},
		{Code: "MX", Name: "Mexico"},
	}

	for i := range countries {
		if err := Upsert(db.WithContext(ctx), &countries[i], "code = ?", countries[i].Code); err != nil {
			return err
		}
	}

	// States for US, CA, AU, IN
	states := []models.State{
		{Name: "California", CountryID: countries[0].ID},
		{Name: "New York", CountryID: countries[0].ID},
		{Name: "Texas", CountryID: countries[0].ID},
		{Name: "Florida", CountryID: countries[0].ID},
		{Name: "Ontario", CountryID: countries[1].ID},
		{Name: "Quebec", CountryID: countries[1].ID},
		{Name: "New South Wales", CountryID: countries[3].ID},
		{Name: "Victoria", CountryID: countries[3].ID},
		{Name: "Maharashtra", CountryID: countries[8].ID},
		{Name: "Karnataka", CountryID: countries[8].ID},
	}

	for i := range states {
		if err := Upsert(db.WithContext(ctx), &states[i], "name = ? AND country_id = ?", states[i].Name, states[i].CountryID); err != nil {
			return err
		}
	}

	// Cities
	cities := []models.City{
		{Name: "San Francisco", StateID: states[0].ID},
		{Name: "Los Angeles", StateID: states[0].ID},
		{Name: "New York City", StateID: states[1].ID},
		{Name: "Buffalo", StateID: states[1].ID},
		{Name: "Houston", StateID: states[2].ID},
		{Name: "Austin", StateID: states[2].ID},
		{Name: "Miami", StateID: states[3].ID},
		{Name: "Orlando", StateID: states[3].ID},
		{Name: "Toronto", StateID: states[4].ID},
		{Name: "Montreal", StateID: states[5].ID},
		{Name: "Sydney", StateID: states[6].ID},
		{Name: "Melbourne", StateID: states[7].ID},
		{Name: "Mumbai", StateID: states[8].ID},
		{Name: "Bangalore", StateID: states[9].ID},
	}

	for i := range cities {
		if err := Upsert(db.WithContext(ctx), &cities[i], "name = ? AND state_id = ?", cities[i].Name, cities[i].StateID); err != nil {
			return err
		}
	}

	return nil
}

// seedUsers creates diverse users with different roles, subscription statuses, and profiles
func seedUsers(db *gorm.DB, ctx context.Context) ([]models.User, error) {
	users := []models.User{
		// Active users with different roles
		{Email: "john.doe@example.com", Phone: "15551230001", PasswordHash: "password", FirstName: "John", LastName: "Doe", IsActive: true, Role: models.RoleUser, Local: "en", EmailVerifiedAt: timePtr(time.Now().AddDate(0, 0, -30))},
		{Email: "jane.smith@example.com", Phone: "15551230002", PasswordHash: "password", FirstName: "Jane", LastName: "Smith", IsActive: true, Role: models.RoleUser, Local: "en", EmailVerifiedAt: timePtr(time.Now().AddDate(0, 0, -15))},
		{Email: "bob.wilson@example.com", Phone: "15551230003", PasswordHash: "password", FirstName: "Bob", LastName: "Wilson", IsActive: true, Role: models.RoleUser, Local: "en", EmailVerifiedAt: timePtr(time.Now().AddDate(0, 0, -7))},
		{Email: "alice.brown@example.com", Phone: "15551230004", PasswordHash: "password", FirstName: "Alice", LastName: "Brown", IsActive: true, Role: models.RoleUser, Local: "en", EmailVerifiedAt: timePtr(time.Now().AddDate(0, 0, -45))},
		{Email: "charlie.davis@example.com", Phone: "15551230005", PasswordHash: "password", FirstName: "Charlie", LastName: "Davis", IsActive: true, Role: models.RoleUser, Local: "en", EmailVerifiedAt: timePtr(time.Now().AddDate(0, 0, -20))},

		// Admin and moderator users
		{Email: "admin@example.com", Phone: "15551230006", PasswordHash: "password", FirstName: "Ada", LastName: "Admin", IsActive: true, Role: models.RoleAdmin, Local: "en", EmailVerifiedAt: timePtr(time.Now().AddDate(0, 0, -60))},
		{Email: "moderator@example.com", Phone: "15551230007", PasswordHash: "password", FirstName: "Mike", LastName: "Moderator", IsActive: true, Role: models.RoleModerator, Local: "en", EmailVerifiedAt: timePtr(time.Now().AddDate(0, 0, -40))},

		// Inactive users
		{Email: "inactive@example.com", Phone: "15551230008", PasswordHash: "password", FirstName: "Inactive", LastName: "User", IsActive: false, Role: models.RoleUser, Local: "en"},

		// Users with unverified emails
		{Email: "unverified@example.com", Phone: "15551230009", PasswordHash: "password", FirstName: "Unverified", LastName: "User", IsActive: true, Role: models.RoleUser, Local: "en"},

		// Users with different locales
		{Email: "french.user@example.com", Phone: "15551230010", PasswordHash: "password", FirstName: "Pierre", LastName: "Dupont", IsActive: true, Role: models.RoleUser, Local: "fr", EmailVerifiedAt: timePtr(time.Now().AddDate(0, 0, -10))},
		{Email: "german.user@example.com", Phone: "15551230011", PasswordHash: "password", FirstName: "Hans", LastName: "Mueller", IsActive: true, Role: models.RoleUser, Local: "de", EmailVerifiedAt: timePtr(time.Now().AddDate(0, 0, -25))},
	}

	for i := range users {
		if err := Upsert(db.WithContext(ctx), &users[i], "email = ?", users[i].Email); err != nil {
			return nil, err
		}
	}

	// Create user profiles
	profiles := []models.UserProfile{
		{UserID: users[0].ID, DisplayName: "JohnDoe", Bio: "Book lover and collector", AvatarURL: "https://example.com/john.jpg", Linkedin: "https://linkedin.com/in/johndoe", CountryID: &[]uint{1}[0], StateID: &[]uint{1}[0], CityID: &[]uint{1}[0]},
		{UserID: users[1].ID, DisplayName: "JaneSmith", Bio: "Sci-fi enthusiast", AvatarURL: "https://example.com/jane.jpg", Linkedin: "https://linkedin.com/in/janesmith", CountryID: &[]uint{1}[0], StateID: &[]uint{2}[0], CityID: &[]uint{3}[0]},
		{UserID: users[2].ID, DisplayName: "BobWilson", Bio: "Mystery novel fan", AvatarURL: "https://example.com/bob.jpg", Linkedin: "https://linkedin.com/in/bobwilson", CountryID: &[]uint{1}[0], StateID: &[]uint{3}[0], CityID: &[]uint{5}[0]},
		{UserID: users[3].ID, DisplayName: "AliceBrown", Bio: "Romance reader", AvatarURL: "https://example.com/alice.jpg", Linkedin: "https://linkedin.com/in/alicebrown", CountryID: &[]uint{1}[0], StateID: &[]uint{4}[0], CityID: &[]uint{7}[0]},
		{UserID: users[4].ID, DisplayName: "CharlieDavis", Bio: "Non-fiction reader", AvatarURL: "https://example.com/charlie.jpg", Linkedin: "https://linkedin.com/in/charliedavis", CountryID: &[]uint{2}[0], StateID: &[]uint{5}[0], CityID: &[]uint{9}[0]},
		{UserID: users[5].ID, DisplayName: "AdaAdmin", Bio: "Site administrator", AvatarURL: "https://example.com/ada.jpg", Linkedin: "https://linkedin.com/in/ada", CountryID: &[]uint{1}[0], StateID: &[]uint{1}[0], CityID: &[]uint{3}[0]},
		{UserID: users[6].ID, DisplayName: "MikeMod", Bio: "Community moderator", AvatarURL: "https://example.com/mike.jpg", Linkedin: "https://linkedin.com/in/mike", CountryID: &[]uint{1}[0], StateID: &[]uint{1}[0], CityID: &[]uint{3}[0]},
		{UserID: users[7].ID, DisplayName: "InactiveUser", Bio: "Former user", AvatarURL: "https://example.com/inactive.jpg", Linkedin: "https://linkedin.com/in/inactive", CountryID: &[]uint{1}[0], StateID: &[]uint{1}[0], CityID: &[]uint{1}[0]},
		{UserID: users[8].ID, DisplayName: "UnverifiedUser", Bio: "New user", AvatarURL: "https://example.com/unverified.jpg", Linkedin: "https://linkedin.com/in/unverified", CountryID: &[]uint{1}[0], StateID: &[]uint{1}[0], CityID: &[]uint{1}[0]},
		{UserID: users[9].ID, DisplayName: "PierreDupont", Bio: "Lecteur français", AvatarURL: "https://example.com/pierre.jpg", Linkedin: "https://linkedin.com/in/pierre", CountryID: &[]uint{6}[0]},
		{UserID: users[10].ID, DisplayName: "HansMueller", Bio: "Deutscher Leser", AvatarURL: "https://example.com/hans.jpg", Linkedin: "https://linkedin.com/in/hans", CountryID: &[]uint{5}[0]},
	}

	for i := range profiles {
		if err := Upsert(db.WithContext(ctx), &profiles[i], "user_id = ?", profiles[i].UserID); err != nil {
			return nil, err
		}
	}

	return users, nil
}

// seedAuthorsAndGenres creates extensive authors and genres for book relationships
func seedAuthorsAndGenres(db *gorm.DB, ctx context.Context) ([]models.Author, []models.Genre, error) {
	// Authors
	authors := []models.Author{
		{Name: "Isaac Asimov"},
		{Name: "George Orwell"},
		{Name: "Mary Shelley"},
		{Name: "J.K. Rowling"},
		{Name: "Stephen King"},
		{Name: "Agatha Christie"},
		{Name: "Jane Austen"},
		{Name: "Charles Dickens"},
		{Name: "Mark Twain"},
		{Name: "Ernest Hemingway"},
		{Name: "F. Scott Fitzgerald"},
		{Name: "Harper Lee"},
		{Name: "Toni Morrison"},
		{Name: "Gabriel García Márquez"},
		{Name: "Milan Kundera"},
		{Name: "Umberto Eco"},
		{Name: "Salman Rushdie"},
		{Name: "Margaret Atwood"},
		{Name: "Neil Gaiman"},
		{Name: "Terry Pratchett"},
	}

	for i := range authors {
		if err := Upsert(db.WithContext(ctx), &authors[i], "name = ?", authors[i].Name); err != nil {
			return nil, nil, err
		}
	}

	// Genres
	genres := []models.Genre{
		{Slug: "fiction", Name: "Fiction"},
		{Slug: "non-fiction", Name: "Non-Fiction"},
		{Slug: "sci-fi", Name: "Science Fiction"},
		{Slug: "fantasy", Name: "Fantasy"},
		{Slug: "mystery", Name: "Mystery"},
		{Slug: "romance", Name: "Romance"},
		{Slug: "thriller", Name: "Thriller"},
		{Slug: "horror", Name: "Horror"},
		{Slug: "biography", Name: "Biography"},
		{Slug: "history", Name: "History"},
		{Slug: "philosophy", Name: "Philosophy"},
		{Slug: "poetry", Name: "Poetry"},
		{Slug: "drama", Name: "Drama"},
		{Slug: "comedy", Name: "Comedy"},
		{Slug: "adventure", Name: "Adventure"},
		{Slug: "young-adult", Name: "Young Adult"},
		{Slug: "children", Name: "Children's"},
		{Slug: "self-help", Name: "Self Help"},
		{Slug: "business", Name: "Business"},
		{Slug: "technology", Name: "Technology"},
	}

	for i := range genres {
		if err := Upsert(db.WithContext(ctx), &genres[i], "slug = ?", genres[i].Slug); err != nil {
			return nil, nil, err
		}
	}

	return authors, genres, nil
}

// Helper function to create time pointers
func timePtr(t time.Time) *time.Time {
	return &t
}

// seedBooks creates books with varied conditions, locations, and relationships
func seedBooks(db *gorm.DB, ctx context.Context, users []models.User, authors []models.Author, genres []models.Genre) ([]models.Book, error) {
	now := time.Now()
	books := []models.Book{
		// Books by different authors with various conditions
		{OwnerID: users[0].ID, Title: "1984", AuthorID: &authors[1].ID, Language: "EN", Condition: "like_new", AvailableFrom: &now, Active: true, LocationCountry: stringPtr("US"), LocationState: stringPtr("California"), LocationCity: stringPtr("San Francisco")},
		{OwnerID: users[1].ID, Title: "Foundation", AuthorID: &authors[0].ID, Language: "EN", Condition: "good", AvailableFrom: &now, Active: true, LocationCountry: stringPtr("US"), LocationState: stringPtr("New York"), LocationCity: stringPtr("New York City")},
		{OwnerID: users[2].ID, Title: "Frankenstein", AuthorID: &authors[2].ID, Language: "EN", Condition: "acceptable", AvailableFrom: &now, Active: true, LocationCountry: stringPtr("US"), LocationState: stringPtr("Texas"), LocationCity: stringPtr("Houston")},
		{OwnerID: users[0].ID, Title: "Harry Potter and the Philosopher's Stone", AuthorID: &authors[3].ID, Language: "EN", Condition: "new", AvailableFrom: &now, Active: true, LocationCountry: stringPtr("US"), LocationState: stringPtr("California"), LocationCity: stringPtr("Los Angeles")},
		{OwnerID: users[1].ID, Title: "The Shining", AuthorID: &authors[4].ID, Language: "EN", Condition: "like_new", AvailableFrom: &now, Active: true, LocationCountry: stringPtr("US"), LocationState: stringPtr("Florida"), LocationCity: stringPtr("Miami")},
		{OwnerID: users[3].ID, Title: "Murder on the Orient Express", AuthorID: &authors[5].ID, Language: "EN", Condition: "good", AvailableFrom: &now, Active: true, LocationCountry: stringPtr("CA"), LocationState: stringPtr("Ontario"), LocationCity: stringPtr("Toronto")},
		{OwnerID: users[4].ID, Title: "Pride and Prejudice", AuthorID: &authors[6].ID, Language: "EN", Condition: "like_new", AvailableFrom: &now, Active: true, LocationCountry: stringPtr("AU"), LocationState: stringPtr("New South Wales"), LocationCity: stringPtr("Sydney")},
		{OwnerID: users[0].ID, Title: "Great Expectations", AuthorID: &authors[7].ID, Language: "EN", Condition: "acceptable", AvailableFrom: &now, Active: true, LocationCountry: stringPtr("US"), LocationState: stringPtr("California"), LocationCity: stringPtr("San Francisco")},
		{OwnerID: users[2].ID, Title: "The Adventures of Tom Sawyer", AuthorID: &authors[8].ID, Language: "EN", Condition: "good", AvailableFrom: &now, Active: true, LocationCountry: stringPtr("US"), LocationState: stringPtr("Texas"), LocationCity: stringPtr("Austin")},
		{OwnerID: users[1].ID, Title: "The Old Man and the Sea", AuthorID: &authors[9].ID, Language: "EN", Condition: "like_new", AvailableFrom: &now, Active: true, LocationCountry: stringPtr("US"), LocationState: stringPtr("New York"), LocationCity: stringPtr("Buffalo")},
		{OwnerID: users[3].ID, Title: "The Great Gatsby", AuthorID: &authors[10].ID, Language: "EN", Condition: "new", AvailableFrom: &now, Active: true, LocationCountry: stringPtr("US"), LocationState: stringPtr("Florida"), LocationCity: stringPtr("Orlando")},
		{OwnerID: users[4].ID, Title: "To Kill a Mockingbird", AuthorID: &authors[11].ID, Language: "EN", Condition: "good", AvailableFrom: &now, Active: true, LocationCountry: stringPtr("AU"), LocationState: stringPtr("Victoria"), LocationCity: stringPtr("Melbourne")},
		{OwnerID: users[0].ID, Title: "Beloved", AuthorID: &authors[12].ID, Language: "EN", Condition: "like_new", AvailableFrom: &now, Active: true, LocationCountry: stringPtr("US"), LocationState: stringPtr("California"), LocationCity: stringPtr("San Francisco")},
		{OwnerID: users[2].ID, Title: "One Hundred Years of Solitude", AuthorID: &authors[13].ID, Language: "EN", Condition: "good", AvailableFrom: &now, Active: true, LocationCountry: stringPtr("US"), LocationState: stringPtr("Texas"), LocationCity: stringPtr("Houston")},
		{OwnerID: users[1].ID, Title: "The Unbearable Lightness of Being", AuthorID: &authors[14].ID, Language: "EN", Condition: "acceptable", AvailableFrom: &now, Active: true, LocationCountry: stringPtr("US"), LocationState: stringPtr("New York"), LocationCity: stringPtr("New York City")},

		// Some archived books for soft delete testing
		{OwnerID: users[0].ID, Title: "Archived Book 1", AuthorID: &authors[0].ID, Language: "EN", Condition: "good", AvailableFrom: &now, Active: false, ArchivedAt: timePtr(now.AddDate(0, 0, -10)), LocationCountry: stringPtr("US"), LocationState: stringPtr("California"), LocationCity: stringPtr("San Francisco")},
		{OwnerID: users[1].ID, Title: "Archived Book 2", AuthorID: &authors[1].ID, Language: "EN", Condition: "like_new", AvailableFrom: &now, Active: false, ArchivedAt: timePtr(now.AddDate(0, 0, -5)), LocationCountry: stringPtr("US"), LocationState: stringPtr("New York"), LocationCity: stringPtr("New York City")},
	}

	for i := range books {
		if err := Upsert(db.WithContext(ctx), &books[i], "owner_id = ? AND title = ?", books[i].OwnerID, books[i].Title); err != nil {
			return nil, err
		}
	}

	// Create book images for each book
	for _, book := range books {
		img := models.BookImage{
			BookID:     book.ID,
			URL:        fmt.Sprintf("https://img.example.com/%d/cover.jpg", book.ID),
			Width:      800,
			Height:     1200,
			IsPrimary:  true,
			UploadedAt: time.Now(),
		}
		_ = Upsert(db.WithContext(ctx), &img, "book_id = ? AND is_primary = 1", book.ID)

		// Add secondary images for some books
		if book.ID%3 == 0 {
			secondaryImg := models.BookImage{
				BookID:     book.ID,
				URL:        fmt.Sprintf("https://img.example.com/%d/back.jpg", book.ID),
				Width:      800,
				Height:     1200,
				IsPrimary:  false,
				UploadedAt: time.Now(),
			}
			_ = Upsert(db.WithContext(ctx), &secondaryImg, "book_id = ? AND is_primary = 0", book.ID)
		}
	}

	// Assign genres to books (many-to-many relationships)
	bookGenreAssignments := [][]int{
		{0, 2},    // 1984 -> Fiction, Sci-Fi
		{1, 2},    // Foundation -> Non-Fiction, Sci-Fi
		{2, 0},    // Frankenstein -> Fiction
		{3, 3},    // Harry Potter -> Fantasy
		{4, 7},    // The Shining -> Horror
		{5, 4},    // Murder on the Orient Express -> Mystery
		{6, 0, 5}, // Pride and Prejudice -> Fiction, Romance
		{7, 0},    // Great Expectations -> Fiction
		{8, 0},    // Tom Sawyer -> Fiction
		{9, 0},    // Old Man and the Sea -> Fiction
		{10, 0},   // Great Gatsby -> Fiction
		{11, 0},   // To Kill a Mockingbird -> Fiction
		{12, 0},   // Beloved -> Fiction
		{13, 0},   // One Hundred Years -> Fiction
		{14, 0},   // Unbearable Lightness -> Fiction
	}

	for i, book := range books[:len(bookGenreAssignments)] {
		var bookGenres []models.Genre
		for _, genreIndex := range bookGenreAssignments[i] {
			if genreIndex < len(genres) {
				bookGenres = append(bookGenres, genres[genreIndex])
			}
		}
		if len(bookGenres) > 0 {
			_ = db.Model(&book).Association("Genres").Replace(bookGenres)
		}
	}

	return books, nil
}

// seedSubscriptions creates multiple subscription plans and various subscription statuses
func seedSubscriptions(db *gorm.DB, ctx context.Context, users []models.User) ([]models.Subscription, error) {
	// Create subscription plans
	plans := []models.SubscriptionPlan{
		{Slug: "free", Name: "Free", PriceCents: 0, Currency: "USD", Interval: "month", Active: true},
		{Slug: "basic", Name: "Basic", PriceCents: 999, Currency: "USD", Interval: "month", Active: true},
		{Slug: "premium", Name: "Premium", PriceCents: 1999, Currency: "USD", Interval: "month", Active: true},
		{Slug: "enterprise", Name: "Enterprise", PriceCents: 4999, Currency: "USD", Interval: "month", Active: true},
		{Slug: "annual-basic", Name: "Basic Annual", PriceCents: 9999, Currency: "USD", Interval: "year", Active: true},
		{Slug: "annual-premium", Name: "Premium Annual", PriceCents: 19999, Currency: "USD", Interval: "year", Active: true},
		{Slug: "inactive-plan", Name: "Inactive Plan", PriceCents: 999, Currency: "USD", Interval: "month", Active: false},
	}

	for i := range plans {
		if err := Upsert(db.WithContext(ctx), &plans[i], "slug = ?", plans[i].Slug); err != nil {
			return nil, err
		}
	}

	// Create subscriptions with different statuses
	now := time.Now()
	subscriptions := []models.Subscription{
		{UserID: uint64(users[0].ID), PlanID: &[]uint64{uint64(plans[1].ID)}[0], Status: models.SubscriptionStatusActive, CurrentPeriodStart: timePtr(now.AddDate(0, 0, -15)), CurrentPeriodEnd: timePtr(now.AddDate(0, 0, 15))},
		{UserID: uint64(users[1].ID), PlanID: &[]uint64{uint64(plans[2].ID)}[0], Status: models.SubscriptionStatusActive, CurrentPeriodStart: timePtr(now.AddDate(0, 0, -10)), CurrentPeriodEnd: timePtr(now.AddDate(0, 0, 20))},
		{UserID: uint64(users[2].ID), PlanID: &[]uint64{uint64(plans[0].ID)}[0], Status: models.SubscriptionStatusTrialing, CurrentPeriodStart: timePtr(now.AddDate(0, 0, -5)), CurrentPeriodEnd: timePtr(now.AddDate(0, 0, 25))},
		{UserID: uint64(users[3].ID), PlanID: &[]uint64{uint64(plans[1].ID)}[0], Status: models.SubscriptionStatusPastDue, CurrentPeriodStart: timePtr(now.AddDate(0, 0, -30)), CurrentPeriodEnd: timePtr(now.AddDate(0, 0, -5))},
		{UserID: uint64(users[4].ID), PlanID: &[]uint64{uint64(plans[2].ID)}[0], Status: models.SubscriptionStatusCanceled, CurrentPeriodStart: timePtr(now.AddDate(0, 0, -60)), CurrentPeriodEnd: timePtr(now.AddDate(0, 0, -30)), CancelAtPeriodEnd: true},
		{UserID: uint64(users[5].ID), PlanID: &[]uint64{uint64(plans[3].ID)}[0], Status: models.SubscriptionStatusActive, CurrentPeriodStart: timePtr(now.AddDate(0, 0, -20)), CurrentPeriodEnd: timePtr(now.AddDate(0, 0, 10))},
		{UserID: uint64(users[6].ID), PlanID: &[]uint64{uint64(plans[1].ID)}[0], Status: models.SubscriptionStatusExpired, CurrentPeriodStart: timePtr(now.AddDate(0, 0, -90)), CurrentPeriodEnd: timePtr(now.AddDate(0, 0, -60))},
		{UserID: uint64(users[7].ID), PlanID: &[]uint64{uint64(plans[0].ID)}[0], Status: models.SubscriptionStatusActive, CurrentPeriodStart: timePtr(now.AddDate(0, 0, -45)), CurrentPeriodEnd: timePtr(now.AddDate(0, 0, 15))},
	}

	for i := range subscriptions {
		if err := Upsert(db.WithContext(ctx), &subscriptions[i], "user_id = ?", subscriptions[i].UserID); err != nil {
			return nil, err
		}
	}

	return subscriptions, nil
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}

// seedCommunities creates multiple communities with different settings
func seedCommunities(db *gorm.DB, ctx context.Context, users []models.User) ([]models.Community, error) {
	communities := []models.Community{
		{Name: "Book Lovers", Slug: "book-lovers", Description: "A place for book lovers to discuss their favorite reads", CreatorID: users[0].ID, RequirePaidChat: false},
		{Name: "Sci-Fi Enthusiasts", Slug: "sci-fi-enthusiasts", Description: "Science fiction book discussions and recommendations", CreatorID: users[1].ID, RequirePaidChat: true},
		{Name: "Mystery Readers", Slug: "mystery-readers", Description: "Mystery and thriller book club", CreatorID: users[2].ID, RequirePaidChat: false},
		{Name: "Romance Book Club", Slug: "romance-book-club", Description: "Romance novel discussions", CreatorID: users[3].ID, RequirePaidChat: true},
		{Name: "Non-Fiction Readers", Slug: "non-fiction-readers", Description: "Non-fiction book discussions", CreatorID: users[4].ID, RequirePaidChat: false},
		{Name: "Classic Literature", Slug: "classic-literature", Description: "Classic literature appreciation society", CreatorID: users[5].ID, RequirePaidChat: true},
		{Name: "Young Adult Books", Slug: "young-adult-books", Description: "YA book discussions", CreatorID: users[6].ID, RequirePaidChat: false},
		{Name: "Local Book Exchange", Slug: "local-book-exchange", Description: "Local book trading community", CreatorID: users[0].ID, RequirePaidChat: false},
	}

	for i := range communities {
		if err := Upsert(db.WithContext(ctx), &communities[i], "slug = ?", communities[i].Slug); err != nil {
			return nil, err
		}
	}

	// Create community memberships
	memberships := []models.CommunityMember{
		{CommunityID: communities[0].ID, UserID: users[0].ID},
		{CommunityID: communities[0].ID, UserID: users[1].ID},
		{CommunityID: communities[0].ID, UserID: users[2].ID},
		{CommunityID: communities[1].ID, UserID: users[1].ID},
		{CommunityID: communities[1].ID, UserID: users[4].ID},
		{CommunityID: communities[2].ID, UserID: users[2].ID},
		{CommunityID: communities[2].ID, UserID: users[3].ID},
		{CommunityID: communities[3].ID, UserID: users[3].ID},
		{CommunityID: communities[3].ID, UserID: users[0].ID},
		{CommunityID: communities[4].ID, UserID: users[4].ID},
		{CommunityID: communities[4].ID, UserID: users[1].ID},
		{CommunityID: communities[5].ID, UserID: users[5].ID},
		{CommunityID: communities[5].ID, UserID: users[6].ID},
		{CommunityID: communities[6].ID, UserID: users[6].ID},
		{CommunityID: communities[6].ID, UserID: users[2].ID},
		{CommunityID: communities[7].ID, UserID: users[0].ID},
		{CommunityID: communities[7].ID, UserID: users[1].ID},
		{CommunityID: communities[7].ID, UserID: users[2].ID},
		{CommunityID: communities[7].ID, UserID: users[3].ID},
	}

	for i := range memberships {
		if err := Upsert(db.WithContext(ctx), &memberships[i], "community_id = ? AND user_id = ?", memberships[i].CommunityID, memberships[i].UserID); err != nil {
			return nil, err
		}
	}

	return communities, nil
}

// seedExchanges creates exchanges in various states for testing
func seedExchanges(db *gorm.DB, ctx context.Context, users []models.User, books []models.Book) ([]models.Exchange, error) {
	now := time.Now()
	exchanges := []models.Exchange{
		// Requested exchanges
		{RequesterID: users[0].ID, ResponderID: &users[1].ID, RequesterBookID: &books[0].ID, ResponderBookID: &books[1].ID, ShippingPayerUserID: users[0].ID, Status: string(models.StatusRequested), RequestedAt: now.AddDate(0, 0, -5), StatusUpdatedAt: now.AddDate(0, 0, -5), Metadata: "{}"},
		{RequesterID: users[2].ID, ResponderID: &users[3].ID, RequesterBookID: &books[2].ID, ResponderBookID: &books[5].ID, ShippingPayerUserID: users[2].ID, Status: string(models.StatusRequested), RequestedAt: now.AddDate(0, 0, -3), StatusUpdatedAt: now.AddDate(0, 0, -3), Metadata: "{}"},

		// Accepted exchanges
		{RequesterID: users[1].ID, ResponderID: &users[4].ID, RequesterBookID: &books[3].ID, ResponderBookID: &books[6].ID, ShippingPayerUserID: users[1].ID, Status: string(models.StatusAccepted), RequestedAt: now.AddDate(0, 0, -10), StatusUpdatedAt: now.AddDate(0, 0, -8), AgreedStartDate: now.AddDate(0, 0, -8), AgreedEndDate: now.AddDate(0, 0, 7), Metadata: "{}"},

		// Shipped exchanges
		{RequesterID: users[0].ID, ResponderID: &users[2].ID, RequesterBookID: &books[7].ID, ResponderBookID: &books[8].ID, ShippingPayerUserID: users[0].ID, Status: string(models.StatusShipped), RequestedAt: now.AddDate(0, 0, -15), StatusUpdatedAt: now.AddDate(0, 0, -2), AgreedStartDate: now.AddDate(0, 0, -12), AgreedEndDate: now.AddDate(0, 0, 3), ShippingProvider: "UPS", ShippingTrackingNumber: "1Z999AA1234567890", ShippingCostCents: 1299, Metadata: "{}"},

		// Completed exchanges
		{RequesterID: users[3].ID, ResponderID: &users[1].ID, RequesterBookID: &books[9].ID, ResponderBookID: &books[10].ID, ShippingPayerUserID: users[3].ID, Status: string(models.StatusCompleted), RequestedAt: now.AddDate(0, 0, -30), StatusUpdatedAt: now.AddDate(0, 0, -5), AgreedStartDate: now.AddDate(0, 0, -25), AgreedEndDate: now.AddDate(0, 0, -5), CompletedAt: now.AddDate(0, 0, -5), Metadata: "{}"},
		{RequesterID: users[4].ID, ResponderID: &users[0].ID, RequesterBookID: &books[11].ID, ResponderBookID: &books[12].ID, ShippingPayerUserID: users[4].ID, Status: string(models.StatusCompleted), RequestedAt: now.AddDate(0, 0, -45), StatusUpdatedAt: now.AddDate(0, 0, -20), AgreedStartDate: now.AddDate(0, 0, -40), AgreedEndDate: now.AddDate(0, 0, -20), CompletedAt: now.AddDate(0, 0, -20), Metadata: "{}"},

		// Canceled exchanges
		{RequesterID: users[2].ID, ResponderID: &users[4].ID, RequesterBookID: &books[13].ID, ResponderBookID: &books[14].ID, ShippingPayerUserID: users[2].ID, Status: string(models.StatusCanceled), RequestedAt: now.AddDate(0, 0, -20), StatusUpdatedAt: now.AddDate(0, 0, -15), CanceledAt: now.AddDate(0, 0, -15), Metadata: "{}"},

		// Disputed exchanges
		{RequesterID: users[1].ID, ResponderID: &users[3].ID, RequesterBookID: &books[1].ID, ResponderBookID: &books[2].ID, ShippingPayerUserID: users[1].ID, Status: string(models.StatusDisputed), RequestedAt: now.AddDate(0, 0, -25), StatusUpdatedAt: now.AddDate(0, 0, -10), AgreedStartDate: now.AddDate(0, 0, -20), AgreedEndDate: now.AddDate(0, 0, -10), DisputeReason: "Book condition not as described", DisputeOpenedAt: now.AddDate(0, 0, -10), Metadata: "{}"},
	}

	for i := range exchanges {
		if err := Upsert(db.WithContext(ctx), &exchanges[i], "requester_id = ? AND responder_id = ? AND status = ?", exchanges[i].RequesterID, exchanges[i].ResponderID, exchanges[i].Status); err != nil {
			return nil, err
		}
	}

	return exchanges, nil
}

// seedMessages creates chat threads and messages for communication testing
func seedMessages(db *gorm.DB, ctx context.Context, users []models.User, exchanges []models.Exchange, communities []models.Community) error {
	// Create chat threads for exchanges
	for i, exchange := range exchanges {
		chat := models.ChatThread{
			ExchangeID: exchange.ID,
			CreatedBy:  exchange.RequesterID,
			Archived:   i%3 == 0, // Some threads are archived
		}
		if err := Upsert(db.WithContext(ctx), &chat, "exchange_id = ?", chat.ExchangeID); err != nil {
			return err
		}

		// Create messages in chat threads
		messages := []models.Message{
			{ThreadID: chat.ID, SenderID: exchange.RequesterID, Type: models.MessageTypeText, Body: "Hi! I'm interested in trading this book.", Attachments: "[]"},
		}

		// Add responder message if responder exists
		if exchange.ResponderID != nil {
			messages = append(messages, models.Message{
				ThreadID:    chat.ID,
				SenderID:    *exchange.ResponderID,
				Type:        models.MessageTypeText,
				Body:        "Sounds good! What's the condition like?",
				Attachments: "[]",
			})
			messages = append(messages, models.Message{
				ThreadID:    chat.ID,
				SenderID:    exchange.RequesterID,
				Type:        models.MessageTypeText,
				Body:        "It's in excellent condition, barely read.",
				Attachments: "[]",
			})
		}

		for _, msg := range messages {
			if err := Upsert(db.WithContext(ctx), &msg, "thread_id = ? AND sender_id = ? AND body = ?", msg.ThreadID, msg.SenderID, msg.Body); err != nil {
				return err
			}
		}
	}

	// Create community threads and messages
	communityThreads := []models.CommunityThread{
		{CommunityID: communities[0].ID, CreatedBy: users[0].ID, Title: "Welcome to Book Lovers!"},
		{CommunityID: communities[0].ID, CreatedBy: users[1].ID, Title: "What are you reading this week?"},
		{CommunityID: communities[1].ID, CreatedBy: users[1].ID, Title: "Best Sci-Fi books of 2023"},
		{CommunityID: communities[2].ID, CreatedBy: users[2].ID, Title: "Mystery recommendations"},
		{CommunityID: communities[3].ID, CreatedBy: users[3].ID, Title: "Romance novel discussions"},
	}

	for i := range communityThreads {
		if err := Upsert(db.WithContext(ctx), &communityThreads[i], "community_id = ? AND title = ?", communityThreads[i].CommunityID, communityThreads[i].Title); err != nil {
			return err
		}
	}

	// Create community messages
	communityMessages := []models.CommunityMessage{
		{ThreadID: communityThreads[0].ID, SenderID: users[0].ID, Body: "Welcome everyone to our book community!"},
		{ThreadID: communityThreads[0].ID, SenderID: users[1].ID, Body: "Thanks for creating this space!"},
		{ThreadID: communityThreads[1].ID, SenderID: users[1].ID, Body: "I'm currently reading 'Dune' - amazing world-building!"},
		{ThreadID: communityThreads[1].ID, SenderID: users[2].ID, Body: "I just finished 'The Martian' - highly recommend!"},
		{ThreadID: communityThreads[2].ID, SenderID: users[1].ID, Body: "Foundation series by Asimov is a must-read!"},
		{ThreadID: communityThreads[3].ID, SenderID: users[2].ID, Body: "Agatha Christie's Poirot series is fantastic!"},
		{ThreadID: communityThreads[4].ID, SenderID: users[3].ID, Body: "Jane Austen's works are timeless classics!"},
	}

	for i := range communityMessages {
		if err := Upsert(db.WithContext(ctx), &communityMessages[i], "thread_id = ? AND sender_id = ? AND body = ?", communityMessages[i].ThreadID, communityMessages[i].SenderID, communityMessages[i].Body); err != nil {
			return err
		}
	}

	return nil
}

// seedReviewsAndRatings creates book reviews and user ratings for analytics
func seedReviewsAndRatings(db *gorm.DB, ctx context.Context, users []models.User, books []models.Book, exchanges []models.Exchange) error {
	// Create book reviews
	reviews := []models.BookReview{
		{BookID: books[0].ID, ReviewerID: users[1].ID, Rating: 5, Comment: "A timeless classic that everyone should read!"},
		{BookID: books[0].ID, ReviewerID: users[2].ID, Rating: 4, Comment: "Thought-provoking and well-written."},
		{BookID: books[1].ID, ReviewerID: users[0].ID, Rating: 5, Comment: "Brilliant science fiction series!"},
		{BookID: books[2].ID, ReviewerID: users[3].ID, Rating: 3, Comment: "Interesting but a bit dated."},
		{BookID: books[3].ID, ReviewerID: users[4].ID, Rating: 5, Comment: "Magical world-building at its finest!"},
		{BookID: books[4].ID, ReviewerID: users[0].ID, Rating: 4, Comment: "Scary and atmospheric."},
		{BookID: books[5].ID, ReviewerID: users[1].ID, Rating: 5, Comment: "Perfect mystery with great characters."},
		{BookID: books[6].ID, ReviewerID: users[2].ID, Rating: 4, Comment: "Classic romance with strong characters."},
		{BookID: books[7].ID, ReviewerID: users[3].ID, Rating: 3, Comment: "Good but quite long."},
		{BookID: books[8].ID, ReviewerID: users[4].ID, Rating: 4, Comment: "Fun adventure story."},
		{BookID: books[9].ID, ReviewerID: users[0].ID, Rating: 5, Comment: "Hemingway at his best!"},
		{BookID: books[10].ID, ReviewerID: users[1].ID, Rating: 4, Comment: "Great American novel."},
		{BookID: books[11].ID, ReviewerID: users[2].ID, Rating: 5, Comment: "Powerful and moving story."},
		{BookID: books[12].ID, ReviewerID: users[3].ID, Rating: 4, Comment: "Beautifully written."},
		{BookID: books[13].ID, ReviewerID: users[4].ID, Rating: 5, Comment: "Magical realism at its finest!"},
	}

	for i := range reviews {
		if err := Upsert(db.WithContext(ctx), &reviews[i], "book_id = ? AND reviewer_id = ?", reviews[i].BookID, reviews[i].ReviewerID); err != nil {
			return err
		}
	}

	// Create user ratings for completed exchanges
	ratings := []models.UserRating{
		{ExchangeID: exchanges[3].ID, RaterID: exchanges[3].RequesterID, RatedUserID: *exchanges[3].ResponderID, Rating: 5, Comment: "Great trade, book arrived in perfect condition!"},
		{ExchangeID: exchanges[3].ID, RaterID: *exchanges[3].ResponderID, RatedUserID: exchanges[3].RequesterID, Rating: 4, Comment: "Smooth transaction, would trade again."},
		{ExchangeID: exchanges[4].ID, RaterID: exchanges[4].RequesterID, RatedUserID: *exchanges[4].ResponderID, Rating: 5, Comment: "Excellent communication and fast shipping!"},
		{ExchangeID: exchanges[4].ID, RaterID: *exchanges[4].ResponderID, RatedUserID: exchanges[4].RequesterID, Rating: 5, Comment: "Perfect trade, highly recommend!"},
	}

	for i := range ratings {
		if err := Upsert(db.WithContext(ctx), &ratings[i], "exchange_id = ? AND rater_id = ? AND rated_user_id = ?", ratings[i].ExchangeID, ratings[i].RaterID, ratings[i].RatedUserID); err != nil {
			return err
		}
	}

	return nil
}

// seedNotifications creates various notification types
func seedNotifications(db *gorm.DB, ctx context.Context, users []models.User) error {
	notifications := []models.Notification{
		{UserID: users[1].ID, Type: models.NotificationTypeExchangeRequestReceived, Payload: `{"message":"You have a new exchange request","exchange_id":1}`, Read: false},
		{UserID: users[1].ID, Type: models.NotificationTypeExchangeRequestAccepted, Payload: `{"message":"Your exchange request was accepted","exchange_id":2}`, Read: true},
		{UserID: users[2].ID, Type: models.NotificationTypeExchangeShipped, Payload: `{"message":"Your book has been shipped","exchange_id":3}`, Read: false},
		{UserID: users[3].ID, Type: models.NotificationTypeExchangeDelivered, Payload: `{"message":"Your book has been delivered","exchange_id":4}`, Read: true},
		{UserID: users[4].ID, Type: models.NotificationTypeExchangeCompleted, Payload: `{"message":"Exchange completed successfully","exchange_id":5}`, Read: false},
		{UserID: users[0].ID, Type: models.NotificationTypeNewMessageInExchange, Payload: `{"message":"You have a new message","exchange_id":1}`, Read: false},
		{UserID: users[1].ID, Type: models.NotificationTypeNewCommunityMessage, Payload: `{"message":"New message in Book Lovers community","community_id":1}`, Read: true},
		{UserID: users[2].ID, Type: models.NotificationTypeBookReviewReceived, Payload: `{"message":"Someone reviewed your book","book_id":1}`, Read: false},
		{UserID: users[3].ID, Type: models.NotificationTypeSubscriptionExpiringSoon, Payload: `{"message":"Your subscription expires in 3 days","subscription_id":1}`, Read: false},
		{UserID: users[4].ID, Type: models.NotificationTypeSubscriptionRenewed, Payload: `{"message":"Your subscription has been renewed","subscription_id":2}`, Read: true},
		{UserID: users[0].ID, Type: models.NotificationTypeGeneralAnnouncement, Payload: `{"message":"New features are now available!","announcement_id":1}`, Read: false},
	}

	for i := range notifications {
		if err := Upsert(db.WithContext(ctx), &notifications[i], "user_id = ? AND type = ?", notifications[i].UserID, notifications[i].Type); err != nil {
			return err
		}
	}

	return nil
}

// seedModeration creates reports and moderation actions
func seedModeration(db *gorm.DB, ctx context.Context, users []models.User) error {
	// Create reports
	reports := []models.Report{
		{ReporterID: users[0].ID, TargetType: "message", TargetID: 1, Reason: "spam", Metadata: "{}", HandledBy: users[5].ID},
		{ReporterID: users[1].ID, TargetType: "user", TargetID: users[2].ID, Reason: "inappropriate_behavior", Metadata: "{}", HandledBy: users[6].ID},
		{ReporterID: users[2].ID, TargetType: "book", TargetID: 1, Reason: "misleading_description", Metadata: "{}", HandledBy: users[5].ID},
		{ReporterID: users[3].ID, TargetType: "exchange", TargetID: 1, Reason: "fraud", Metadata: "{}", HandledBy: users[6].ID},
		{ReporterID: users[4].ID, TargetType: "community_message", TargetID: 1, Reason: "harassment", Metadata: "{}", HandledBy: users[5].ID},
	}

	for i := range reports {
		if err := Upsert(db.WithContext(ctx), &reports[i], "reporter_id = ? AND target_type = ? AND target_id = ?", reports[i].ReporterID, reports[i].TargetType, reports[i].TargetID); err != nil {
			return err
		}
	}

	// Create moderation actions
	moderationActions := []models.ModerationAction{
		{TargetType: "message", TargetID: 1, Action: "review", Reason: "routine", Metadata: "{}"},
		{TargetType: "user", TargetID: users[2].ID, Action: "warn", Reason: "inappropriate_behavior", Metadata: "{}"},
		{TargetType: "book", TargetID: 1, Action: "flag", Reason: "misleading_description", Metadata: "{}"},
		{TargetType: "exchange", TargetID: 1, Action: "investigate", Reason: "fraud", Metadata: "{}"},
		{TargetType: "community_message", TargetID: 1, Action: "remove", Reason: "harassment", Metadata: "{}"},
	}

	for i := range moderationActions {
		if err := Upsert(db.WithContext(ctx), &moderationActions[i], "target_type = ? AND target_id = ?", moderationActions[i].TargetType, moderationActions[i].TargetID); err != nil {
			return err
		}
	}

	return nil
}

// seedActivityLogs creates activity logs for audit trails
func seedActivityLogs(db *gorm.DB, ctx context.Context, users []models.User, books []models.Book) error {
	activityLogs := []models.ActivityLog{
		{UserID: &users[0].ID, Action: models.ActionCreate, ObjectType: "book", ObjectID: &books[0].ID, Payload: `{"title":"1984"}`},
		{UserID: &users[1].ID, Action: models.ActionUpdate, ObjectType: "user_profile", ObjectID: &users[1].ID, Payload: `{"bio":"Updated bio"}`},
		{UserID: &users[2].ID, Action: models.ActionDelete, ObjectType: "book", ObjectID: &books[15].ID, Payload: `{"title":"Archived Book 1"}`},
		{UserID: &users[3].ID, Action: models.ActionCreate, ObjectType: "exchange", ObjectID: &[]uint{1}[0], Payload: `{"exchange_id":1}`},
		{UserID: &users[4].ID, Action: models.ActionUpdate, ObjectType: "subscription", ObjectID: &[]uint{1}[0], Payload: `{"status":"active"}`},
		{UserID: &users[0].ID, Action: models.ActionCreate, ObjectType: "community", ObjectID: &[]uint{1}[0], Payload: `{"name":"Book Lovers"}`},
		{UserID: &users[1].ID, Action: models.ActionUpdate, ObjectType: "book", ObjectID: &books[1].ID, Payload: `{"title":"Foundation"}`},
		{UserID: &users[2].ID, Action: models.ActionDelete, ObjectType: "message", ObjectID: &[]uint{1}[0], Payload: `{"message_id":1}`},
	}

	for i := range activityLogs {
		if err := Upsert(db.WithContext(ctx), &activityLogs[i], "action = ? AND object_type = ? AND object_id = ?", activityLogs[i].Action, activityLogs[i].ObjectType, activityLogs[i].ObjectID); err != nil {
			return err
		}
	}

	return nil
}

// seedMessageQuotaUsage creates message quota usage records
func seedMessageQuotaUsage(db *gorm.DB, ctx context.Context, users []models.User) error {
	now := time.Now()
	quotaUsages := []models.MessageQuotaUsage{
		{UserID: users[0].ID, PeriodStart: now.AddDate(0, 0, -30), PeriodEnd: now.AddDate(0, 0, -1), MessagesSent: 45},
		{UserID: users[1].ID, PeriodStart: now.AddDate(0, 0, -30), PeriodEnd: now.AddDate(0, 0, -1), MessagesSent: 32},
		{UserID: users[2].ID, PeriodStart: now.AddDate(0, 0, -30), PeriodEnd: now.AddDate(0, 0, -1), MessagesSent: 28},
		{UserID: users[3].ID, PeriodStart: now.AddDate(0, 0, -30), PeriodEnd: now.AddDate(0, 0, -1), MessagesSent: 15},
		{UserID: users[4].ID, PeriodStart: now.AddDate(0, 0, -30), PeriodEnd: now.AddDate(0, 0, -1), MessagesSent: 38},
		{UserID: users[5].ID, PeriodStart: now.AddDate(0, 0, -30), PeriodEnd: now.AddDate(0, 0, -1), MessagesSent: 12},
		{UserID: users[6].ID, PeriodStart: now.AddDate(0, 0, -30), PeriodEnd: now.AddDate(0, 0, -1), MessagesSent: 8},
		{UserID: users[0].ID, PeriodStart: now.AddDate(0, 0, -60), PeriodEnd: now.AddDate(0, 0, -31), MessagesSent: 52},
		{UserID: users[1].ID, PeriodStart: now.AddDate(0, 0, -60), PeriodEnd: now.AddDate(0, 0, -31), MessagesSent: 41},
		{UserID: users[2].ID, PeriodStart: now.AddDate(0, 0, -60), PeriodEnd: now.AddDate(0, 0, -31), MessagesSent: 35},
	}

	for i := range quotaUsages {
		if err := Upsert(db.WithContext(ctx), &quotaUsages[i], "user_id = ? AND period_start = ?", quotaUsages[i].UserID, quotaUsages[i].PeriodStart); err != nil {
			return err
		}
	}

	return nil
}

func Upsert[T any](db *gorm.DB, model T, where string, args ...any) error {
	return db.Where(where, args...).Assign(model).FirstOrCreate(model).Error
}
