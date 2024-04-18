package models

import (
	"gorm.io/gorm"
	// "time"
)

type Table struct {
	gorm.Model
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

// Override the default table name.
func (Table) TableName() string {
	return "tables"
}