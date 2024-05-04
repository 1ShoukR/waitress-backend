package models

import (
	"gorm.io/gorm"
	"time"
)

// Receipt represents a receipt record in the database.
type Receipt struct {
    gorm.Model
    TipAmount      *float64        // Pointer to allow nil (nullable)
    AssignedWaiter uint            `gorm:"not null"`
    AssignedUser   uint            `gorm:"not null"`
    RestaurantID   uint            `gorm:"not null"`
    Restaurant     Restaurant      `gorm:"foreignKey:RestaurantID"`
}

// Restaurant represents a restaurant record in the database.
type Restaurant struct {
    gorm.Model
    OwnerID         uint            `gorm:"not null"`
    Name            string          `gorm:"size:255;not null"`
    Address         string          `gorm:"size:255;not null"`
    Phone           string          `gorm:"size:255;not null"`
    Email           string          `gorm:"size:255;not null"`
    Website         *string         // Pointer to allow nil (nullable)
    NumberOfTables  *int            // Pointer to allow nil (nullable)
    Latitude        *float64        // Pointer to allow nil (nullable)
    Longitude       *float64        // Pointer to allow nil (nullable)
    Receipts        []Receipt       `gorm:"foreignKey:RestaurantID"` // One-to-many relationship
    Reservations    *[]Reservation  `gorm:"foreignKey:RestaurantID"`
    Owner           User            `gorm:"foreignKey:OwnerID"`
}

// Reservation represents a reservation record in the database.
type Reservation struct {
    gorm.Model
    ReservationID           uint            `gorm:"primaryKey"`
    RestaurantID            uint            `gorm:"not null"`
    UserID                  uint            `gorm:"not null"` // This field links the reservation to a user
    TableID                 uint            `gorm:"not null"`
    Time                    time.Time       `gorm:"not null"`
	Restaurant     			Restaurant      `gorm:"foreignKey:RestaurantID"`

    // User                    User            `gorm:"foreignKey:UserID"` // GORM uses this for loading the user
}

// MenuItem represents a menu item record in the database.
type MenuItem struct {
	gorm.Model
	MenuID         uint            `gorm:"primaryKey"`
	RestaurantID   uint            `gorm:"not null"`
	NameOfItem     *string         // Pointer to allow nil (nullable)
	Price          *float64        // Pointer to allow nil (nullable)
	IsAvailable    bool            `gorm:"default:true"`
	Restaurant     Restaurant      `gorm:"foreignKey:RestaurantID"`
}

// Order represents an order record in the database.
type Order struct {
	gorm.Model
	OrderID        uint            `gorm:"primaryKey"`
	ReservationID  uint            `gorm:"not null"`
	UserID         uint            `gorm:"not null"`
	Total          *float64        // Pointer to allow nil (nullable)
	IsPaid         bool            `gorm:"default:false"`
	Reservation    Reservation     `gorm:"foreignKey:ReservationID"`
}