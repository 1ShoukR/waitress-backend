// This file contains models related to a restaurant and its related entities
//
// The models here are as follows:
// - Receipt
// - Restaurant
// - Reservation
// - MenuItem
// - Order
// - Table

package models

import (
	"time"

	"gorm.io/gorm"
)

// Receipt represents a receipt record in the database.
type Receipt struct {
	gorm.Model
	TipAmount      *float64   // Pointer to allow nil (nullable)
	AssignedWaiter uint       `gorm:"not null"`
	AssignedUser   uint       `gorm:"not null"`
	RestaurantID   uint       `gorm:"not null"`
	Restaurant     Restaurant `gorm:"foreignKey:RestaurantID"`
}

// Category represents a category in which a restaurant can be classified.
type Category struct {
	CategoryID   uint    `gorm:"primaryKey;autoIncrement:true"`
	CategoryName string  `gorm:"size:255;not null"`
	ImageURL     *string // Pointer to allow nil (nullable)
}

// Restaurant represents a restaurant record in the database.
type Restaurant struct {
	RestaurantId   uint           `gorm:"primaryKey;autoIncrement:true"`
	OwnerID        uint           `gorm:"index"`
	Name           string         `gorm:"size:255;not null"`
	Address        string         `gorm:"size:255;not null"`
	Phone          string         `gorm:"size:255;not null"`
	Email          string         `gorm:"size:255;not null"`
	Website        *string        // Pointer to allow nil (nullable)
	Categories     []Category     `gorm:"many2many:restaurant_categories;"`
	NumberOfTables *int           // Pointer to allow nil (nullable)
	Latitude       *float64       // Pointer to allow nil (nullable)
	Longitude      *float64       // Pointer to allow nil (nullable)
	Receipts       []Receipt      `gorm:"foreignKey:RestaurantID"` // One-to-many relationship
	Reservations   *[]Reservation `gorm:"foreignKey:RestaurantID"`
	MenuItems      *[]MenuItem    `gorm:"foreignKey:RestaurantID"` // One-to-many relationship
	Owner          User           `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Ratings        *[]Rating      `gorm:"foreignKey:RestaurantID"` // One-to-many relationship
	ImageURL       *string        // Pointer to allow nil (nullable)
	// Calculated fields
	AverageRating float32 `gorm:"default:0"`
	ReviewCount   *int    `gorm:"-"`
}

func (Restaurant) TableName() string {
	return "restaurants"
}

// Rating represents a rating record in the database.
// We can use this to calculate the average rating for a restaurant.
type Rating struct {
	RatingID     uint       `gorm:"primaryKey;autoIncrement:true"`
	Comment      string     `gorm:"size:255"`
	Rating       uint       `gorm:"not null"`
	RestaurantID uint       `gorm:"not null"`
	UserID       uint       `gorm:"not null"`
	Restaurant   Restaurant `gorm:"foreignKey:RestaurantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// Reservation represents a reservation record in the database.
type Reservation struct {
	ReservationID uint           `gorm:"primaryKey;autoIncrement:true"`
	CreatedAt     time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	RestaurantID  uint           `gorm:"not null"`
	UserID        uint           `gorm:"not null"`
	TableID       uint           `gorm:"not null"`
	Time          time.Time      `gorm:"not null"`
	Restaurant    Restaurant     `gorm:"foreignKey:RestaurantID"`
	// User            User                     `gorm:"foreignKey:UserID"`
}

// MenuItem represents a menu item record in the database.
type MenuItem struct {
	MenuID       uint       `gorm:"primaryKey;autoIncrement:true"`
	RestaurantID uint       `gorm:"not null"`
	NameOfItem   *string    // Pointer to allow nil (nullable)
	Price        *float64   // Pointer to allow nil (nullable)
	IsAvailable  bool       `gorm:"default:true"`
	Category     *string    // Pointer to allow nil (nullable)
	ImageURL     *string    // Pointer to allow nil (nullable)
	Description  *string    // Pointer to allow nil (nullable)
	Restaurant   Restaurant `gorm:"foreignKey:RestaurantID"`
}

func (m *MenuItem) GetMenuItem(db *gorm.DB, menuId string) (*MenuItem, error) {
	var menuItem MenuItem
	if err := db.Where("menu_id = ?", menuId).First(&menuItem).Error; err != nil {
		return nil, err
	}
	return &menuItem, nil
}

// Order represents an order record in the database.
type Order struct {
	OrderID       uint     `gorm:"primaryKey;autoIncrement:true"`
	ReservationID uint     `gorm:"not null"`
	UserID        uint     `gorm:"not null"`
	Total         *float64 // Pointer to allow nil (nullable)
	IsPaid        bool     `gorm:"default:false"`
	// Reservation    Reservation     `gorm:"foreignKey:ReservationID"`
}

// Table represents a table record in the database.
type Table struct {
	TableID             uint   `gorm:"primaryKey;autoIncrement"`
	RestaurantID        uint   `gorm:"not null"`
	ReservationID       *uint  // It's a pointer to allow nil value when no reservation is associated
	TableNumber         uint   `gorm:"not null"`
	Capacity            uint   `gorm:"not null"`
	LocationDescription string `gorm:"size:200"` // Description of the table's location
	IsReserved          bool   `gorm:"default:false"`
	CustomerID          *uint  // It's a pointer to allow nil value when no customer is associated

	// Define relationships
	// Restaurant   Restaurant `gorm:"foreignKey:RestaurantID"`
	// Reservation  Reservation `gorm:"foreignKey:ReservationID"`
	// Customer     User        `gorm:"foreignKey:CustomerID"`
}

// CalcAvgRating calculates the average rating for a restaurant.
func (r *Restaurant) CalcAvgRating(db *gorm.DB, restaurantId string) (float32, error) {
	var avgRating float32
	err := db.Table("rating").Select("AVG(rating) as average_rating").Where("restaurant_id = ?", restaurantId).Row().Scan(&avgRating)
	if err != nil {
		return 0, err
	}
	return avgRating, nil
}

// UpdateAvgRating updates the average rating for a restaurant.
func (r *Restaurant) UpdateAvgRating(db *gorm.DB, restaurantId string, avgRating float32) error {
	if err := db.Model(&Restaurant{}).Where("restaurant_id = ?", restaurantId).Update("average_rating", avgRating).Error; err != nil {
		return err
	}
	return nil
}
