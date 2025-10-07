package seeder

import (
	"context"
	"fmt"
	"time"

	"github.com/Amanuel-0/gorm-pg/internals/database/models"
	"gorm.io/gorm"
)

// SeedAll seeds the database with relational data. It is idempotent.
func SeedAll(db *gorm.DB) error {
	ctx := context.Background()

	// locations
	country := models.Country{Code: "US", Name: "United States"}
	if err := Upsert(db.WithContext(ctx), &country, "code = ?", country.Code); err != nil {
		return err
	}
	state := models.State{Name: "California", CountryID: country.ID}
	if err := Upsert(db.WithContext(ctx), &state, "name = ? AND country_id = ?", state.Name, country.ID); err != nil {
		return err
	}
	city := models.City{Name: "San Francisco", StateID: state.ID}
	if err := Upsert(db.WithContext(ctx), &city, "name = ? AND state_id = ?", city.Name, state.ID); err != nil {
		return err
	}

	// users and profiles
	users := []models.User{
		{Email: "doe@example.com", Phone: "15551230001", PasswordHash: "password", FirstName: "John", LastName: "Doe", IsActive: true, Role: models.RoleUser, Local: "en"},
		{Email: "smith@example.com", Phone: "15551230002", PasswordHash: "password", FirstName: "Jane", LastName: "Smith", IsActive: true, Role: models.RoleUser, Local: "en"},
		{Email: "admin@example.com", Phone: "15551230003", PasswordHash: "password", FirstName: "Ada", LastName: "Admin", IsActive: true, Role: models.RoleAdmin, Local: "en"},
	}
	for i := range users {
		if err := Upsert(db.WithContext(ctx), &users[i], "email = ?", users[i].Email); err != nil {
			return err
		}
	}

	profiles := []models.UserProfile{
		{UserID: users[0].ID, Bio: "Bio of John Doe", AvatarURL: "https://example.com/john.jpg", Linkedin: "https://linkedin.com/in/johndoe", CountryID: &country.ID, StateID: &state.ID, CityID: &city.ID},
		{UserID: users[1].ID, Bio: "Bio of Jane Smith", AvatarURL: "https://example.com/jane.jpg", Linkedin: "https://linkedin.com/in/janesmith", CountryID: &country.ID, StateID: &state.ID, CityID: &city.ID},
		{UserID: users[2].ID, Bio: "Site administrator", AvatarURL: "https://example.com/ada.jpg", Linkedin: "https://linkedin.com/in/ada", CountryID: &country.ID, StateID: &state.ID, CityID: &city.ID},
	}
	for i := range profiles {
		if err := Upsert(db.WithContext(ctx), &profiles[i], "user_id = ?", profiles[i].UserID); err != nil {
			return err
		}
	}

	// genres and authors
	genres := []models.Genre{{Slug: "fiction", Name: "Fiction"}, {Slug: "non-fiction", Name: "Non-Fiction"}, {Slug: "sci-fi", Name: "Sci-Fi"}}
	for i := range genres {
		if err := Upsert(db.WithContext(ctx), &genres[i], "slug = ?", genres[i].Slug); err != nil {
			return err
		}
	}

	authors := []models.Author{{Name: "Isaac Asimov"}, {Name: "George Orwell"}, {Name: "Mary Shelley"}}
	for i := range authors {
		if err := Upsert(db.WithContext(ctx), &authors[i], "name = ?", authors[i].Name); err != nil {
			return err
		}
	}

	// books
	now := time.Now()
	books := []models.Book{
		{OwnerID: users[0].ID, Title: "1984", AuthorID: &authors[1].Model.ID, Language: "EN", AvailableFrom: &now, Active: true},
		{OwnerID: users[1].ID, Title: "Foundation", AuthorID: &authors[0].Model.ID, Language: "EN", AvailableFrom: &now, Active: true},
		{OwnerID: users[0].ID, Title: "Frankenstein", AuthorID: &authors[2].Model.ID, Language: "EN", AvailableFrom: &now, Active: true},
	}
	for i := range books {
		if err := Upsert(db.WithContext(ctx), &books[i], "owner_id = ? AND title = ?", books[i].OwnerID, books[i].Title); err != nil {
			return err
		}
	}

	// book genres m2m
	_ = db.Model(&books[0]).Association("Genres").Replace([]models.Genre{genres[0]})
	_ = db.Model(&books[1]).Association("Genres").Replace([]models.Genre{genres[2]})
	_ = db.Model(&books[2]).Association("Genres").Replace([]models.Genre{genres[0], genres[2]})

	// book images
	for _, b := range books {
		img := models.BookImage{BookID: b.ID, URL: fmt.Sprintf("https://img.example.com/%d/cover.jpg", b.ID), Width: 800, Height: 1200, IsPrimary: true, UploadedAt: time.Now()}
		_ = Upsert(db.WithContext(ctx), &img, "book_id = ? AND is_primary = 1", b.ID)
	}

	// book reviews
	review := models.BookReview{BookID: books[0].ID, ReviewerID: users[1].ID, Rating: 5, Comment: "A timeless classic"}
	if err := Upsert(db.WithContext(ctx), &review, "book_id = ? AND reviewer_id = ?", review.BookID, review.ReviewerID); err != nil {
		return err
	}

	// subscription plans and subscriptions
	basicPlan := models.SubscriptionPlan{Slug: "basic", Name: "Basic", PriceCents: 0, Currency: "USD", Interval: "month", Active: true}
	if err := Upsert(db.WithContext(ctx), &basicPlan, "slug = ?", basicPlan.Slug); err != nil {
		return err
	}

	planID64 := uint64(basicPlan.ID)
	sub := models.Subscription{UserID: uint64(users[0].ID), PlanID: &planID64, Status: models.SubscriptionStatusActive}
	if err := Upsert(db.WithContext(ctx), &sub, "user_id = ?", sub.UserID); err != nil {
		return err
	}

	// community, threads, messages
	community := models.Community{Name: "Book Lovers", Slug: "book-lovers", Description: "A place for book lovers", CreatorID: users[0].ID}
	if err := Upsert(db.WithContext(ctx), &community, "slug = ?", community.Slug); err != nil {
		return err
	}
	// ensure membership
	member := models.CommunityMember{CommunityID: community.ID, UserID: users[0].ID}
	if err := Upsert(db.WithContext(ctx), &member, "community_id = ? AND user_id = ?", member.CommunityID, member.UserID); err != nil {
		return err
	}

	thread := models.CommunityThread{CommunityID: community.ID, CreatedBy: users[0].ID, Title: "Welcome"}
	if err := Upsert(db.WithContext(ctx), &thread, "community_id = ? AND title = ?", thread.CommunityID, thread.Title); err != nil {
		return err
	}

	cm := models.CommunityMessage{ThreadID: thread.ID, SenderID: users[0].ID, Body: "Hello everyone!"}
	if err := Upsert(db.WithContext(ctx), &cm, "thread_id = ? AND sender_id = ?", cm.ThreadID, cm.SenderID); err != nil {
		return err
	}

	// exchange + chat thread + message
	exch := models.Exchange{RequesterID: users[0].ID, ResponderID: &users[1].ID, RequesterBookID: &books[0].ID, ResponderBookID: &books[1].ID, ShippingPayerUserID: users[0].ID, Status: string(models.StatusRequested), Metadata: "{}"}
	if err := Upsert(db.WithContext(ctx), &exch, "requester_id = ? AND responder_id = ? AND status = ?", exch.RequesterID, exch.ResponderID, exch.Status); err != nil {
		return err
	}

	chat := models.ChatThread{ExchangeID: exch.ID, CreatedBy: users[0].ID, Archived: false}
	if err := Upsert(db.WithContext(ctx), &chat, "exchange_id = ?", chat.ExchangeID); err != nil {
		return err
	}

	msg := models.Message{ThreadID: chat.ID, SenderID: users[0].ID, Type: models.MessageTypeText, Body: "Interested in trading?", Attachments: "[]"}
	if err := Upsert(db.WithContext(ctx), &msg, "thread_id = ? AND sender_id = ?", msg.ThreadID, msg.SenderID); err != nil {
		return err
	}

	// ratings
	rating := models.UserRating{RaterID: users[0].ID, RatedUserID: users[1].ID, ExchangeID: exch.ID, Rating: 5, Comment: "Great trade!"}
	if err := Upsert(db.WithContext(ctx), &rating, "rater_id = ? AND rated_user_id = ? AND exchange_id = ?", rating.RaterID, rating.RatedUserID, rating.ExchangeID); err != nil {
		return err
	}

	// notifications
	notif := models.Notification{UserID: users[1].ID, Type: models.NotificationTypeNewMessageInExchange, Payload: "{\"message\":\"You have a new message\"}", Read: false}
	if err := Upsert(db.WithContext(ctx), &notif, "user_id = ? AND type = ?", notif.UserID, notif.Type); err != nil {
		return err
	}

	// moderation action
	mod := models.ModerationAction{TargetType: "message", TargetID: msg.ID, Action: "review", Reason: "routine", Metadata: "{}"}
	if err := Upsert(db.WithContext(ctx), &mod, "target_type = ? AND target_id = ?", mod.TargetType, mod.TargetID); err != nil {
		return err
	}

	// report
	rep := models.Report{ReporterID: users[0].ID, TargetType: "message", TargetID: msg.ID, Reason: "spam", Metadata: "{}", HandledBy: users[2].ID}
	if err := Upsert(db.WithContext(ctx), &rep, "reporter_id = ? AND target_type = ? AND target_id = ?", rep.ReporterID, rep.TargetType, rep.TargetID); err != nil {
		return err
	}

	// activity log
	act := models.ActivityLog{UserID: &users[0].ID, Action: models.ActionCreate, ObjectType: "book", ObjectID: &books[0].ID, Payload: "{}"}
	if err := Upsert(db.WithContext(ctx), &act, "action = ? AND object_type = ? AND object_id = ?", act.Action, act.ObjectType, act.ObjectID); err != nil {
		return err
	}

	// message quota usage
	mqu := models.MessageQuotaUsage{UserID: users[0].ID, PeriodStart: time.Now().AddDate(0, 0, -7), PeriodEnd: time.Now(), MessagesSent: 3}
	if err := Upsert(db.WithContext(ctx), &mqu, "user_id = ? AND period_start = ?", mqu.UserID, mqu.PeriodStart); err != nil {
		return err
	}

	return nil
}

func Upsert[T any](db *gorm.DB, model T, where string, args ...any) error {
	return db.Where(where, args...).Assign(model).FirstOrCreate(model).Error
}

// a generic method to handle upsert for the specified model
// currently not using it
func Upsert2[T any](db *gorm.DB, model T, where string, args ...any) error {
	existing := new(T)

	err := db.Where(where, args...).First(existing).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err // actual query error
	}

	if err == nil {
		// Record exists: extract ID and assign to the input model (if it has one)
		db.Model(model).Select("id").Updates(existing)
		return nil
	}

	// No record found, create a new one
	return db.Create(model).Error
}
