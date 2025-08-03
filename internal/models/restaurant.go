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
	TableID       uint   `gorm:"primaryKey;autoIncrement"`
	RestaurantID  uint   `gorm:"not null"`
	ReservationID *uint  // It's a pointer to allow nil value when no reservation is associated
	TableNumber   string `gorm:"not null"` // Changed from uint to string for flexible naming like "T1", "Window-3"
	Capacity      uint   `gorm:"not null"`

	LocationZone        string `gorm:"size:50"`  // "inside", "outside", "patio", "bar"
	LocationDescription string `gorm:"size:200"` // "corner booth", "center dining", "by kitchen"
	ViewDescription     string `gorm:"size:200"` // "street view", "garden view", "no view"
	TableType           string `gorm:"size:50"`  // "booth", "standard", "high-top", "bar-seat"

	IsAvailable bool  `gorm:"default:true"`  // Can table be reserved (not broken, etc.)
	IsReserved  bool  `gorm:"default:false"` // Currently reserved
	CustomerID  *uint // Current customer if occupied

	CoordinateX *float64 `gorm:"default:null"` // Grid X position
	CoordinateY *float64 `gorm:"default:null"` // Grid Y position
	Width       *float64 `gorm:"default:null"` // Table width in grid units
	Height      *float64 `gorm:"default:null"` // Table height in grid units
	Rotation    *float64 `gorm:"default:null"` // Rotation angle in degrees

	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Define relationships
	Restaurant  Restaurant   `gorm:"foreignKey:RestaurantID"`
	Reservation *Reservation `gorm:"foreignKey:ReservationID"`
	Customer    *User        `gorm:"foreignKey:CustomerID"`
}

// RestaurantLayout represents the visual layout configuration for a restaurant
type RestaurantLayout struct {
	LayoutID        uint      `gorm:"primaryKey;autoIncrement"`
	RestaurantID    uint      `gorm:"not null;unique"` // One layout per restaurant
	FloorplanURL    *string   `gorm:"size:500"`        // Future: uploaded floorplan image
	GridWidth       *int      `gorm:"default:null"`    // Grid dimensions for coordinate system
	GridHeight      *int      `gorm:"default:null"`
	ScaleFactor     *float64  `gorm:"default:null"` // Pixels per grid unit
	BackgroundColor *string   `gorm:"size:7"`       // Hex color for background
	CreatedAt       time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`

	// Relationships
	Restaurant Restaurant `gorm:"foreignKey:RestaurantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (RestaurantLayout) TableName() string {
	return "restaurant_layouts"
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
