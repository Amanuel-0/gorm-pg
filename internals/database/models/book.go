package models

import (
	"time"

	"gorm.io/gorm"
)

// create an enum for the condition
type Condition string

const (
	ConditionNew        Condition = "new"
	ConditionLikeNew    Condition = "like_new"
	ConditionGood       Condition = "good"
	ConditionAcceptable Condition = "acceptable"
)

func (c Condition) IsValid() bool {
	switch c {
	case ConditionNew, ConditionLikeNew, ConditionGood, ConditionAcceptable:
		return true
	}
	return false
}

type Book struct {
	// gorm.Model
	ID              uint       `json:"id" gorm:"primaryKey"`
	OwnerID         uint       `json:"owner_id,omitempty" gorm:"not null"`
	Title           string     `json:"title,omitempty" gorm:"type:varchar(1000);not null;uniqueIndex:idx_owner_title"`
	Subtitle        *string    `json:"subtitle,omitempty" gorm:"type:varchar(1000)"`
	AuthorID        *uint      `json:"author_id,omitempty"`
	ISBN            *string    `json:"isbn,omitempty" gorm:"type:varchar(32)"`
	Description     string     `json:"description,omitempty" gorm:"type:text"`
	Language        string     `json:"language,omitempty" gorm:"type:varchar(8);default:'EN'"`
	Condition       Condition  `json:"condition,omitempty" gorm:"type:enum('new','like_new','good','acceptable');default:'good'"`
	AvailableFrom   *time.Time `json:"available_from,omitempty"`
	AvailableUntil  *time.Time `json:"available_until,omitempty"`
	LocationCity    *string    `json:"location_city,omitempty" gorm:"type:varchar(100)"`
	LocationState   *string    `json:"location_state,omitempty" gorm:"type:varchar(100)"`
	LocationCountry *string    `json:"location_country,omitempty" gorm:"type:varchar(100)"`
	Latitude        *float64   `json:"latitude,omitempty"`
	Longitude       *float64   `json:"longitude,omitempty"`
	Location        *string    `json:"location,omitempty" gorm:"type:point"`
	Active          bool       `json:"active,omitempty" gorm:"default:true"`
	ArchivedAt      *time.Time `json:"archived_at,omitempty"`
	PreferredTitles *string    `json:"preferred_titles,omitempty" gorm:"type:json"`

	// timestamps
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	// using `gorm.DeletedAt` makes it possible to soft delete a user
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`

	// Relationships
	Owner       *User         `json:"owner,omitempty" gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE"`
	Author      *Author       `json:"author,omitempty" gorm:"foreignKey:AuthorID;constraint:OnDelete:SET NULL"`
	Genres      []*Genre      `json:"genres,omitempty" gorm:"many2many:book_genres;"`
	Images      []*BookImage  `json:"images,omitempty" gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE"`
	BookReviews []*BookReview `json:"book_reviews,omitempty" gorm:"foreignKey:BookID"`
}

//
//
// I'm going to store the books_count in the User table.
//
//

// // AfterCreate: increment owner's counter
// func (b *Book) AfterCreate(tx *gorm.DB) (err error) {
// 	if b.OwnerID == 0 {
// 		return nil
// 	}
// 	return tx.Model(&User{}).
// 		Where("id = ?", b.OwnerID).
// 		UpdateColumn("books_count", gorm.Expr("books_count + ?", 1)).Error
// }

// // AfterDelete: decrement owner's counter
// func (b *Book) AfterDelete(tx *gorm.DB) (err error) {
// 	if b.OwnerID == 0 {
// 		return nil
// 	}
// 	return tx.Model(&User{}).
// 		Where("id = ?", b.OwnerID).
// 		UpdateColumn("books_count", gorm.Expr("GREATEST(books_count - ?, 0)", 1)).Error
// }

// // BeforeUpdate: record previous owner id into the statement context
// func (b *Book) BeforeUpdate(tx *gorm.DB) (err error) {
// 	var old Book
// 	if err := tx.Unscoped().Select("owner_id").Where("id = ?", b.ID).Take(&old).Error; err == nil {
// 		tx.InstanceSet("old_owner_id", old.OwnerID)
// 	}
// 	return nil
// }

// // AfterUpdate: if owner changed, update both old and new owner's counters
// func (b *Book) AfterUpdate(tx *gorm.DB) (err error) {
// 	v, ok := tx.InstanceGet("old_owner_id")
// 	if !ok {
// 		return nil
// 	}
// 	oldOwnerID, _ := v.(uint)
// 	newOwnerID := b.OwnerID

// 	// If owner didn't change, nothing to do.
// 	if oldOwnerID == newOwnerID {
// 		return nil
// 	}

// 	// Decrement old owner
// 	if oldOwnerID != 0 {
// 		if err := tx.Model(&User{}).
// 			Where("id = ?", oldOwnerID).
// 			UpdateColumn("books_count", gorm.Expr("GREATEST(books_count - ?, 0)", 1)).Error; err != nil {
// 			return err
// 		}
// 	}
// 	// Increment new owner
// 	if newOwnerID != 0 {
// 		if err := tx.Model(&User{}).
// 			Where("id = ?", newOwnerID).
// 			UpdateColumn("books_count", gorm.Expr("books_count + ?", 1)).Error; err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
